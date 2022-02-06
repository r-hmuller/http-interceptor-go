package main

import (
	"crypto/tls"
	"github.com/gorilla/mux"
	"httpInterceptor/checkpoint"
	"httpInterceptor/config"
	"httpInterceptor/handler"
	"httpInterceptor/heartbeat"
	"log"
	"net/http"
	"sync"
)

func main() {
	if config.GetApplicationURL() == "" {
		panic("Couldn't find the HOST ENV")
	}

	if config.GetInterceptorPort() == "" {
		panic("Couldn't find the PORT ENV")
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go startListener()
	if config.GetCheckpointEnabled() {
		wg.Add(1)
		go heartbeat.Monitor()
		wg.Add(1)
		go checkpoint.Monitor()
	}
	wg.Wait()
}

func startListener() {
	// Disable SSL validation, because some client may have invalid certificates
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(handler.InterceptorHandler)
	log.Fatal(http.ListenAndServe(config.GetInterceptorPort(), router))
}
