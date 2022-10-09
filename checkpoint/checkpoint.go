package checkpoint

import (
	"bytes"
	"encoding/json"
	"httpInterceptor/config"
	"httpInterceptor/logging"
	"io"
	"log"
	"net/http"
	"time"
)

func Monitor() {
	period := config.GetSnapshotPeriodicity()
	periodicity := time.Duration(period) * time.Millisecond
	for _ = range time.Tick(periodicity) {
		generateSnapshot()
	}
}

func generateSnapshot() {
	logging.LogToSnapshotFile("Starting snapshot")
	UpdateVariable("isCheckpointing", "true")
	//Mandar requisição state manager
	//POST /service/checkpoint
	sendRequestToStateManager("true")

	postBody, _ := json.Marshal(map[string]string{
		"Namespace":  config.GetContainerNamespace(),
		"Container":  config.GetContainerServiceName(),
		"Service":    config.GetServiceName(),
		"YamlString": config.GetYamlString(),
	})
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	fullUrl := config.GetDaemonEndpoint() + "/snapshot"
	resp, err := http.Post(fullUrl, "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//Mandar requisição state manager
	sendRequestToStateManager("false")

	UpdateVariable("isCheckpointing", "false")
	logging.LogToSnapshotFile("Snapshot completed")
}

func sendRequestToStateManager(status string) {
	postStateManager, _ := json.Marshal(map[string]string{
		"Key":   "",
		"Value": status,
	})
	stateManagerBodyRequest := bytes.NewBuffer(postStateManager)
	stateManagerUrl := config.GetStateManagerUrl() + "/" + config.GetServiceName() + "/checkpoint"
	stateManageResp, err := http.Post(stateManagerUrl, "application/json", stateManagerBodyRequest)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer stateManageResp.Body.Close()
	_, err = io.ReadAll(stateManageResp.Body)
	if err != nil {
		log.Fatal(err)
	}
}
