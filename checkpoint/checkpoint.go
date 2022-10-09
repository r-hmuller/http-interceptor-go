package checkpoint

import (
	"bytes"
	"encoding/json"
	"httpInterceptor/config"
	"httpInterceptor/logging"
	"io"
	"log"
	"net/http"
	"strconv"
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
	sendRequestToStateManager("requestProcessed", strconv.Itoa(GetLatestRequestNumber()))

	//POST /service/checkpoint
	sendRequestToStateManager("checkpoint", "true")

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
		sendRequestToStateManager("checkpoint", "false")
	}

	//Mandar requisição state manager
	sendRequestToStateManager("checkpoint", "false")

	//Aqui marcar todos os processed como snapshoted

	UpdateVariable("isCheckpointing", "false")
	go RemoveAllSnapshotedRequestsFromMaps()
	logging.LogToSnapshotFile("Snapshot completed")
}

//Action = checkpoint or requestsProcessed
func sendRequestToStateManager(action string, value string) {
	postStateManager, _ := json.Marshal(map[string]string{
		"Key":   action,
		"Value": value,
	})
	stateManagerBodyRequest := bytes.NewBuffer(postStateManager)
	stateManagerUrl := config.GetStateManagerUrl() + "/" + config.GetServiceName() + "/config"
	if action == "checkpoint" {
		stateManagerUrl = config.GetStateManagerUrl() + "/" + config.GetServiceName() + "/checkpoint"
	}
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
