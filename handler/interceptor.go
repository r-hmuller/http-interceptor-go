package handler

import (
	"bytes"
	"crypto/tls"
	"errors"
	uuid "github.com/nu7hatch/gouuid"
	"httpInterceptor/config"
	"httpInterceptor/logging"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var requestsMap = make(map[string]*http.Request)
var processedMap = make(map[string]bool)
var someMapMutex = sync.RWMutex{}

type HTTPResponse struct {
	StatusCode         int
	Header             http.Header
	Body               []byte
	InterceptorControl *uuid.UUID
}

func InterceptorHandler(w http.ResponseWriter, r *http.Request) {
	u, err := uuid.NewV4()
	if err != nil {
		logging.LogToFile(err.Error(), "fatal")
	}

	someMapMutex.Lock()
	requestsMap[u.String()] = r
	processedMap[u.String()] = false
	someMapMutex.Unlock()

	requestToApp := r.Clone(r.Context())
	requestToApp.URL.Host = config.GetApplicationURL()
	serverResponse := HTTPResponse{}
	method := requestToApp.Method
	serverResponse = sendRequest(method, requestToApp, u)

	w.WriteHeader(serverResponse.StatusCode)
	_, err = w.Write(serverResponse.Body)
	if err != nil {
		logging.LogToFile(err.Error(), "fatal")
	}

	someMapMutex.Lock()
	processedMap[u.String()] = true
	someMapMutex.Unlock()
}

func sendRequest(method string, destiny *http.Request, uuid *uuid.UUID) HTTPResponse {
	response := HTTPResponse{}
	client := getClient()
	fullUrl := getScheme() + destiny.URL.String()

	requestBody, err := ioutil.ReadAll(destiny.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		response.StatusCode = 500
		return response
	}
	destiny.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
	req, err := http.NewRequest(method, fullUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		response.StatusCode = 500
		return response
	}

	req.Header.Add("Interceptor-Controller", uuid.String())
	addHeaders(destiny, req)
	resp, err := client.Do(req)
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		response.StatusCode = 500
		return response
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logging.LogToFile(err.Error(), "default")
		}
	}(resp.Body)

	body, err := getBodyContent(resp)

	response.StatusCode = resp.StatusCode
	response.Header = resp.Header
	response.Body = body
	response.InterceptorControl = uuid

	return response
}

func getScheme() string {
	scheme := config.GetHttpScheme()
	return scheme
}

func getClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 10000,
		IdleConnTimeout:     5 * time.Second,
		DisableCompression:  true,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}
	return &http.Client{Transport: tr}

}

func getBodyContent(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.LogToFile("Error parsing body", "default")
		return nil, errors.New("Error parsing request body")
	}
	response.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, nil
}

func addHeaders(original *http.Request, created *http.Request) {
	for name, values := range original.Header {
		for _, value := range values {
			created.Header.Add(name, value)
		}
	}
}
