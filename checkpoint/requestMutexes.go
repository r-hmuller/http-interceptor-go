package checkpoint

import (
	"net/http"
	"sync"
)

var requestsMap = make(map[int]*http.Request)
var processedMap = make(map[int]string) //1) Pending 2) Processed 3)Checkpointed
var requestNumber int = 0
var requestsMutex = &sync.RWMutex{}

func SaveRequestToBuffer(request *http.Request) int {
	requestsMutex.Lock()
	requestNumber++
	requestsMap[requestNumber] = request
	processedMap[requestNumber] = "pending"
	requestsMutex.Unlock()
	return requestNumber
}

func UpdateRequestToProcessed(number int) {
	requestsMutex.Lock()
	processedMap[number] = "processed"
	requestsMutex.Unlock()
}

func GetLatestRequestNumber() int {
	return requestNumber
}

func RemoveAllSnapshotedRequestsFromMaps() {
	requestsMutex.Lock()
	items := make([]int, 1000)
	for key, value := range processedMap {
		if value == "snapshoted" {
			_ = append(items, key)
		}
	}
	for _, request := range items {
		delete(requestsMap, request)
		delete(processedMap, request)
	}

	requestsMutex.Unlock()
}

func UpdateRequestToSnapshoted() {
	requestsMutex.Lock()
	for key, value := range processedMap {
		if value == "processed" {
			processedMap[key] = "snapshoted"
		}
	}

	requestsMutex.Unlock()
}
