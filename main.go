package main

import (
	"ablyprojects/deployment"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ably/ably-go/ably"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	URL := os.Getenv("SEELF_URL")
	TOKEN := os.Getenv("SEELF_TOKEN")
	APP_ID := os.Getenv("SEELF_APP_ID")
	ENVIRONMENT := os.Getenv("APP_ENVIRONMENT")
	BRANCH := os.Getenv("APP_BRANCH")
	ABLY := os.Getenv("ABLY_TOKEN")

	// Connect to Ably with your API key
	client, err := ably.NewRealtime(ably.WithKey(ABLY), ably.WithAutoConnect(false))
	if err != nil {
		panic(err)
	}
	client.Connection.OnAll(func(change ably.ConnectionStateChange) {
		fmt.Printf("Connection event: %s state=%s reason=%s\n", change.Event, change.Current, change.Reason)
	})
	client.Connect()

	// Create a channel called 'get-started' and register a listener to subscribe to all messages with the name 'first'
	channel := client.Channels.Get("autodeploy")
	_, err = channel.Subscribe(context.Background(), "helloworld", func(msg *ably.Message) {
		fmt.Printf("Received message : %v\n", msg.Data)

		// Create deployment
		deployResponse, err := deployment.CreateDeployment(URL, TOKEN, APP_ID, ENVIRONMENT, BRANCH, "")
		if err != nil {
			log.Fatalf("Failed to create deployment: %v", err)
		}

		fmt.Printf("Deployment #%d created, waiting for it to complete...\n", deployResponse.DeploymentNumber)

		// Wait for deployment to complete
		err = deployment.WaitForDeployment(URL, TOKEN, APP_ID, deployResponse.DeploymentNumber)
		if err != nil {
			log.Fatalf("Deployment failed: %v", err)
		}
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Subscrition starts")

	for {
	}
}
