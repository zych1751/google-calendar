package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const (
	port = 9000
)

func main() {
	credentialsFile, err := ioutil.ReadFile("./credentials.json")
	if err != nil {
		log.Fatal("Credentials file not found.")
	}
	var googleCredentials GoogleCredentials
	json.Unmarshal(credentialsFile, &googleCredentials)

	googleClient := NewGoogleClient(&googleCredentials)

	// Echo instance
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://schedule.zychspace.com", "http://localhost"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/schedule", getScheduleHandler(googleClient))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
