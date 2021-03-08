package main

import (
	"fmt"
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

// TODO: error handling
func (c GoogleAuthJwtClaims) Valid() error {
	return nil
}

func NewGoogleClient(credentials *GoogleCredentials) *GoogleClient {
	client := &GoogleClient{
		credentials: credentials,
	}
	return client
}

// TODO: make private
func (client *GoogleClient) GetAccessToken() (string, error) {
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
		Scope:     "https://www.googleapis.com/auth/calendar.readonly",
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
