package main

import (
	"github.com/gorilla/mux"
	"httpInterceptor/handler"
	"log"
	"net/http"
	"os"
)

var hostEnv = os.Getenv("HOST")
var portEnv = os.Getenv("PORT")

func main() {
	if hostEnv == "" {
		panic("Couldn't find the HOST ENV")
	}

	if portEnv == "" {
		panic("Couldn't find the PORT ENV")
	}

	router := mux.NewRouter()
	router.PathPrefix("/").HandlerFunc(handler.InterceptorHandler)
	log.Fatal(http.ListenAndServe(portEnv, router))
}
