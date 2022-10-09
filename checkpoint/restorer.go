package checkpoint

import (
	"bytes"
	"encoding/json"
	uuid "github.com/nu7hatch/gouuid"
	"httpInterceptor/config"
	"httpInterceptor/logging"
	"io/ioutil"
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
	reprocessPendingOrProcessedRequests()
	logging.LogToSnapshotFile("Restore completed")

}

func reprocessPendingOrProcessedRequests() {
	reprocessableList := GetReprocessableRequests()
	for _, item := range reprocessableList {
		reprocessItem(item)
	}
}

func reprocessItem(original *http.Request) {
	client := config.GetHttpClient()
	fullUrl := config.GetScheme() + original.URL.String()

	requestBody, err := ioutil.ReadAll(original.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		return
	}
	original.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	req, err := http.NewRequest(original.Method, fullUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		return
	}

	newUuid, _ := uuid.NewV4()
	req.Header.Add("Interceptor-Controller", newUuid.String())
	config.AddHeaders(original, req)
	resp, err := client.Do(req)
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		return
	}

	_, err = config.GetBodyContent(resp)
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		return
	}
	err = resp.Body.Close()
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		return
	}
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		return
	}
}
