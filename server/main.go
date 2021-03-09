package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func main() {
	credentialsFile, err := ioutil.ReadFile("./credentials.json")
	if err != nil {
		log.Fatal("Credentials file not found.")
	}
	var googleCredentials GoogleCredentials
	json.Unmarshal(credentialsFile, &googleCredentials)

	googleClient := NewGoogleClient(&googleCredentials)

	// example
	startTime := time.Now().AddDate(0, 0, -7)
	endTime := time.Now().AddDate(0, 0, 7)

	items, err := googleClient.GetSchedule(startTime, endTime)
	fmt.Println(items)
}
