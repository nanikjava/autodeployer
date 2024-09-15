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

func getEnvVar(key string) string {
    value := os.Getenv(key)
    if value == "" {
        log.Fatalf("Error: Environment variable %s not set", key)
    }
    return value
}

func main() {
	err := godotenv.Load()
	if err != nil {
        log.Println("Error loading .env file, continuing with system environment variables")
	}

	URL := getEnvVar("SEELF_URL")
	TOKEN :=getEnvVar("SEELF_TOKEN")
	APP_ID := getEnvVar("SEELF_APP_ID")
	ENVIRONMENT := getEnvVar("APP_ENVIRONMENT")
	BRANCH := getEnvVar("APP_BRANCH")
	ABLY := getEnvVar("ABLY_TOKEN")

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
