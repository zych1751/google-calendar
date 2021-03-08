package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	credentialsFile, err := ioutil.ReadFile("./credentials.json")
	if err != nil {
		log.Fatal("Credentials file not found.")
	}
	var googleCredentials GoogleCredentials
	json.Unmarshal(credentialsFile, &googleCredentials)
	fmt.Println(googleCredentials)
}
