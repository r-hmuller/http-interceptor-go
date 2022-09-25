package checkpoint

import (
	"bytes"
	"encoding/json"
	"httpInterceptor/config"
	"httpInterceptor/logging"
	"log"
	"net/http"
)

func Restore() {
	logging.LogToSnapshotFile("Starting restore")
	postBody, _ := json.Marshal(map[string]string{
		"Namespace": config.GetContainerNamespace(),
		"Container": config.GetContainerServiceName(),
		"Service":   config.GetServiceName(),
	})
	responseBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request
	fullUrl := config.GetDaemonEndpoint() + "/restore"
	resp, err := http.Post(fullUrl, "application/json", responseBody)
	//Handle Error
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()
	logging.LogToSnapshotFile("Restore completed")
}
