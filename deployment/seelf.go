package deployment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// DeploymentResponse holds the response for deployment creation and status.
type DeploymentResponse struct {
	DeploymentNumber int `json:"deployment_number"`
	State            struct {
		Status    int    `json:"status"`
		ErrorCode string `json:"error_code"`
		Services  []struct {
			Entrypoints []struct {
				IsCustom bool   `json:"is_custom"`
				URL      string `json:"url"`
			} `json:"entrypoints"`
		} `json:"services"`
	} `json:"state"`
}

// CreateDeployment creates a new deployment and returns the deployment number.
func CreateDeployment(url, token, appID, environment, branch, hash string) (*DeploymentResponse, error) {
	payload := map[string]interface{}{
		"environment": environment,
		"git": map[string]interface{}{
			"branch": branch,
		},
	}
	if hash != "" {
		payload["git"].(map[string]interface{})["hash"] = hash
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	deploymentResponse, err := makeRequest("POST", fmt.Sprintf("%s/api/v1/apps/%s/deployments", url, appID), token, payloadBytes)
	if err != nil {
		return nil, err
	}

	if deploymentResponse.DeploymentNumber == 0 {
		return nil, fmt.Errorf("failed to create the deployment")
	}

	return deploymentResponse, nil
}

// CheckDeploymentStatus checks the status of a deployment and returns the updated deployment response.
func CheckDeploymentStatus(url, token, appID string, deploymentNumber int) (*DeploymentResponse, error) {
	return makeRequest("GET", fmt.Sprintf("%s/api/v1/apps/%s/deployments/%d", url, appID, deploymentNumber), token, nil)
}

// WaitForDeployment polls the deployment status until it succeeds or fails.
func WaitForDeployment(url, token, appID string, deploymentNumber int) error {
	for {
		deploymentResponse, err := CheckDeploymentStatus(url, token, appID, deploymentNumber)
		if err != nil {
			return err
		}

		if deploymentResponse.State.Status == 2 {
			return fmt.Errorf("deployment failed: %s", deploymentResponse.State.ErrorCode)
		} else if deploymentResponse.State.Status == 3 {
			log.Println("Deployment succeeded!")
			for _, service := range deploymentResponse.State.Services {
				for _, entrypoint := range service.Entrypoints {
					if !entrypoint.IsCustom {
						fmt.Printf("url=%s\n", entrypoint.URL)
					}
				}
			}
			return nil
		}

		time.Sleep(3 * time.Second)
	}
}

// Helper function to make HTTP requests and process responses.
func makeRequest(method, url, token string, body []byte) (*DeploymentResponse, error) {
	var deploymentResponse DeploymentResponse

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected response code: %d", resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bodyBytes, &deploymentResponse)
	if err != nil {
		return nil, err
	}

	return &deploymentResponse, nil
}
