package handler

import (
	"bytes"
	uuid "github.com/nu7hatch/gouuid"
	"httpInterceptor/checkpoint"
	"httpInterceptor/config"
	"httpInterceptor/logging"
	"io/ioutil"
	"log"
	"net/http"
)

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

	requestNumber := checkpoint.SaveRequestToBuffer(r)

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

	checkpoint.UpdateRequestToProcessed(requestNumber)
}

func sendRequest(method string, destiny *http.Request, uuid *uuid.UUID) HTTPResponse {
	response := HTTPResponse{}
	client := config.GetHttpClient()
	fullUrl := config.GetScheme() + destiny.URL.String()

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
	config.AddHeaders(destiny, req)
	resp, err := client.Do(req)
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		response.StatusCode = 500
		return response
	}
	response.StatusCode = resp.StatusCode
	response.Header = resp.Header
	body, err := config.GetBodyContent(resp)
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		response.StatusCode = 500
		return response
	}
	err = resp.Body.Close()
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		response.StatusCode = 500
		return response
	}

	response.Body = body
	response.InterceptorControl = uuid

	return response
}
