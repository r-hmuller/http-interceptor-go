package handler

import (
	"bytes"
	"errors"
	uuid "github.com/nu7hatch/gouuid"
	"httpInterceptor/logging"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var hostEnv = os.Getenv("HOST")

var requestsMap = make(map[string]*http.Request)
var processedMap = make(map[string]bool)

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
	requestsMap[u.String()] = r
	processedMap[u.String()] = false

	requestToApp := r.Clone(r.Context())
	requestToApp.URL.Host = hostEnv
	serverResponse := HTTPResponse{}
	switch method := requestToApp.Method; method {
	case "GET":
		serverResponse = sendGetRequest(requestToApp, u)
	default:
		panic("Method not found")
	}

	w.WriteHeader(serverResponse.StatusCode)
	_, err = w.Write(serverResponse.Body)
	if err != nil {
		logging.LogToFile(err.Error(), "fatal")
	}

	processedMap[u.String()] = true
}

func sendGetRequest(destiny *http.Request, uuid *uuid.UUID) HTTPResponse {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	client := &http.Client{Transport: tr}
	fullUrl := getScheme(destiny) + destiny.URL.String()
	logging.LogToFile(fullUrl, "default")
	req, err := http.NewRequest("GET", fullUrl, nil)
	req.Header.Add("Interceptor-Controller", uuid.String())
	if err != nil {
		logging.LogToFile(err.Error(), "fatal")
	}

	response := HTTPResponse{}
	resp, err := client.Do(req)
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		response.StatusCode = 500
		return response
	}

	body, err := getBodyContent(resp)

	response.StatusCode = resp.StatusCode
	response.Header = resp.Header
	response.Body = body
	response.InterceptorControl = uuid

	//logging.LogToFile(resp.Status, "default")
	return response
}

func getScheme(destiny *http.Request) string {
	scheme := "https:"
	if destiny.TLS == nil {
		scheme = "http:"
	}
	return scheme
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
