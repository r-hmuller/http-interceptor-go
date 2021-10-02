package main

import (
	"github.com/gorilla/mux"
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
	wg.Add(1)
	go heartbeat.Monitor()
	wg.Wait()
}

func startListener() {
	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(handler.InterceptorHandler)
	log.Fatal(http.ListenAndServe(config.GetInterceptorPort(), router))
}
