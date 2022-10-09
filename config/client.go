package config

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var singleInstance *http.Client

func GetHttpClient() *http.Client {
	tr := &http.Transport{
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 500000,
		IdleConnTimeout:     5 * time.Second,
		DisableCompression:  true,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			singleInstance = &http.Client{Transport: tr}
		}
	}

	return singleInstance

}
