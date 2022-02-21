package checkpoint

import (
	"bytes"
	"encoding/json"
	"httpInterceptor/config"
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
}
