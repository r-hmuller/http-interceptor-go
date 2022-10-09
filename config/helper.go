package config

import (
	"errors"
	"httpInterceptor/logging"
	"io/ioutil"
	"net/http"
)

func GetScheme() string {
	scheme := GetHttpScheme()
	return scheme
}

func GetBodyContent(response *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logging.LogToFile("Error parsing body", "default")
		return nil, errors.New("Error parsing request body")
	}
	return body, nil
}

func AddHeaders(original *http.Request, created *http.Request) {
	for name, values := range original.Header {
		for _, value := range values {
			created.Header.Add(name, value)
		}
	}
}
