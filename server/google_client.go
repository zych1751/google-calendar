package main

import (
	"fmt"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-resty/resty/v2"
)

type GoogleCredentials struct {
	PrivateKeyId string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
}

type GoogleClient struct {
	credentials *GoogleCredentials
}

type GoogleAuthJwtClaims struct {
	Issuer    string `json:"iss,omitempty"`
	Scope     string `json:"scope,omitempty"`
	Audience  string `json:"aud,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
}

type GoogleAuthResponse struct {
	AccessToken string `json:"access_token"`
}

type GoogleCalendarResponse struct {
	Items []GoogleCalendarItem `json:"items"`
}

type GoogleCalendarItem struct {
	Summary string             `json:"summary"`
	Start   GoogleCalendarTime `json:"start"`
	End     GoogleCalendarTime `json:"end"`
}

type GoogleCalendarTime struct {
	DateTime string `json:"dateTime"`
}

// TODO: error handling
func (c GoogleAuthJwtClaims) Valid() error {
	return nil
}

const (
	calendarId      = "zych1751@gmail.com"
	secretRuleRegex = "^!.*$"
	secretText      = "[비밀]"
)

func NewGoogleClient(credentials *GoogleCredentials) *GoogleClient {
	client := &GoogleClient{
		credentials: credentials,
	}
	return client
}

func applySecretRuleRegex(resp *GoogleCalendarResponse) {
	for i := range resp.Items {
		matched, _ := regexp.MatchString(secretRuleRegex, resp.Items[i].Summary)
		if matched {
			resp.Items[i].Summary = secretText
		}
	}
}

// TODO: make private
func (client *GoogleClient) GetSchedule(startTime time.Time, endTime time.Time) ([]GoogleCalendarItem, error) {
	accessToken, err := client.getAccessToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://www.googleapis.com/calendar/v3/calendars/%s/events", calendarId)
	startTimeStr := startTime.Format(time.RFC3339)
	endTimeStr := endTime.Format(time.RFC3339)

	headers := map[string]string{
		"Content-Type":  "application/json; charset=UTF-8",
		"Authorization": "Bearer " + accessToken,
		"X-GFE-SSL":     "yes",
	}
	params := map[string]string{
		"timeMin": startTimeStr,
		"timeMax": endTimeStr,
	}

	httpClient := resty.New()
	resp := GoogleCalendarResponse{}
	_, err = httpClient.
		R().
		SetHeaders(headers).
		SetQueryParams(params).
		SetResult(&resp).
		Get(url)

	applySecretRuleRegex(&resp)
	return resp.Items, nil
}

func (client *GoogleClient) getAccessToken() (string, error) {
	jwt, err := client.getSignedJWT()
	if err != nil {
		return jwt, err
	}

	httpClient := resty.New()
	resp := GoogleAuthResponse{}

	_, err = httpClient.
		R().
		SetResult(&resp).
		SetBody(fmt.Sprintf(`{"grant_type":"urn:ietf:params:oauth:grant-type:jwt-bearer","assertion":"%s"}`, jwt)).
		Post("https://oauth2.googleapis.com/token")
	return resp.AccessToken, err
}

func (client *GoogleClient) getSignedJWT() (string, error) {
	token := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodRS256.Name))
	token.Header["kid"] = client.credentials.PrivateKeyId

	iat := time.Now()
	token.Claims = &GoogleAuthJwtClaims{
		Issuer:    client.credentials.ClientEmail,
		Scope:     "https://www.googleapis.com/auth/calendar.readonly https://www.googleapis.com/auth/calendar https://www.googleapis.com/auth/calendar.events https://www.googleapis.com/auth/calendar.events.readonly",
		Audience:  "https://oauth2.googleapis.com/token",
		IssuedAt:  iat.Unix(),
		ExpiresAt: iat.Add(time.Second * 3600).Unix(),
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(client.credentials.PrivateKey))
	if err != nil {
		return "", err
	}
	return token.SignedString(signKey)
}
