package heartbeat

import (
	"httpInterceptor/checkpoint"
	"httpInterceptor/config"
	"httpInterceptor/logging"
	"net/http"
	"os"
	"strconv"
	"time"
)

var failedTries int = 0

func Monitor() {
	period := getPeriodicity()
	checkServer(time.Duration(period) * time.Millisecond)
}

func checkServer(d time.Duration) {
	for _ = range time.Tick(d) {
		heartbeat()
	}
}

func heartbeat() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", config.GetApplicationURL(), nil)
	if err != nil {
		logging.LogToFile(err.Error(), "default")
	}
	response, err := client.Do(req)
	if err != nil {
		logging.LogToFile(err.Error(), "default")
		failedTries++
		return
	}
	statusOK := response.StatusCode >= 200 && response.StatusCode < 300
	if !statusOK {
		failedTries++
		return
	}
	shouldOpenCircuit()
}

func shouldOpenCircuit() {
	threshold := config.GetHeartbearThreshold()
	if failedTries > threshold {
		checkpoint.Restore()
	}
}

func getPeriodicity() int {
	periodicity := 1000
	aux := os.Getenv("HEARTBEAT_PERIODICITY")
	if aux != "" {
		periodicity, err := strconv.Atoi(aux)
		if err != nil {
			logging.LogToFile(err.Error(), "fatal")
		}
		return periodicity
	}

	return periodicity
}
