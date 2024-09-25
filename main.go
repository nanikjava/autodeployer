package main

import (
	"ablyprojects/deployment"
	"ablyprojects/pkg"
	"ablyprojects/pkg/client"
	"fmt"

	"context"
	"io/ioutil"
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
	d := loadConfigs()
	fmt.Println(d)

	err := godotenv.Load()
	if err != nil {
        log.Println("Error loading .env file, continuing with system environment variables")
	}

	URL := getEnvVar("SEELF_URL")
	TOKEN :=getEnvVar("SEELF_TOKEN")
	ENVIRONMENT := getEnvVar("APP_ENVIRONMENT")
	ABLY := getEnvVar("ABLY_TOKEN")

	c, err  := client.NewRealAblyClient(ABLY)
	if err != nil {
		panic(err)
	}

	ctx:=context.Background()
	err = c.Subscribe(ctx,"autodeploy",func(msg *ably.Message) {
		// msg.name - name
		// msg.data - branch		
		log.Printf("Received message : %v\n", msg)

		p:=availableProject(msg.Name,d)

		if (p.Name=="") {
			log.Printf("Project not found - %s", msg.Name)
			return
		}

		branch, ok := msg.Data.(string); if !ok{
			log.Printf("Failure in converting msg.data %v", err)
			return
		}

		deployResponse, err := deployment.CreateDeployment(URL, TOKEN, p.SeelfKey, ENVIRONMENT, branch, "")
		if err != nil {
			log.Printf("Failed to create deployment: %v", err)
		}

		log.Printf("Deployment #%d created, waiting for it to complete...\n", deployResponse.DeploymentNumber)

		err = deployment.WaitForDeployment(URL, TOKEN,  p.SeelfKey, deployResponse.DeploymentNumber)
		if err != nil {
			log.Printf("Deployment failed: %v", err)
		}
	})

	if err != nil {
		panic(err)
	}
	log.Println("Subscription starts")

	for {
	}
}

func availableProject(name string, d *pkg.Deployment) pkg.Project {
	for _, project:= range d.Deployments.Projects {
		if project.Name == name {
			return project
		}
	}
	return pkg.Project{}
}

func loadConfigs() *pkg.Deployment {
	// Read the file
	fileData, err := ioutil.ReadFile("deployment.yaml")
	if err != nil {
		return nil
	}
	result, err := pkg.ParseConfigs([]byte(fileData))
	if err != nil {
		log.Fatalf("error: %v", err)
		os.Exit(1)
	}

	return result
}
