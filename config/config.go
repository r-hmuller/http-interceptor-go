package config

import (
	"httpInterceptor/logging"
	"os"
	"strconv"
	"strings"
)

var port = "3000"

func GetApplicationURL() string {
	return os.Getenv("HOST") + ":" + port
}

func GetApplicationPort() string {
	if port == "3000" {
		newPort := os.Getenv("APPLICATION_PORT")
		port = newPort
		return newPort
	}
	return port
}

func SetApplicationPort(newPort string) {
	port = newPort
}

func GetInterceptorPort() string {
	return os.Getenv("PORT")
}

func GetHeartbearThreshold() int {
	threshold := 30
	thresholdEnv := os.Getenv("HEARTBEAT_THRESHOLD")
	if thresholdEnv != "" {
		threshold, err := strconv.Atoi(thresholdEnv)
		if err != nil {
			logging.LogToFile(err.Error(), "fatal")
		}
		return threshold
	}
	return threshold
}

func GetSnapshotPeriodicity() int {
	periodicity := 1800
	periodicityEnv := os.Getenv("SNAPSHOT_PERIODICITY")
	if periodicityEnv != "" {
		periodicity, err := strconv.Atoi(periodicityEnv)
		if err != nil {
			logging.LogToFile(err.Error(), "fatal")
		}
		return periodicity
	}
	return periodicity
}

func GetPodName() string {
	return os.Getenv("POD_NAME")
}

func GetRegistry() string {
	return os.Getenv("REGISTRY")
}

func GetCheckpointEnabled() bool {
	shouldDoCheckpoint := true
	doCheckpoint := os.Getenv("ENABLE_CHECKPOINT")
	if doCheckpoint != "" {
		shouldDoCheckpoint, err := strconv.ParseBool(doCheckpoint)
		if err != nil {
			panic(err)
		}
		return shouldDoCheckpoint
	}
	return shouldDoCheckpoint
}

func GetStateManagerUrl() string {
	return os.Getenv("STATE_MANAGER")
}

func GetLogginPath() string {
	return os.Getenv("LOGGING_PATH")
}

func GetHttpScheme() string {
	scheme := "http:"
	schemeEnv := os.Getenv("HTTP_SCHEME")
	if schemeEnv == "https" {
		return "https:"
	}
	return scheme
}

func GetContainerServiceName() string {
	return os.Getenv("CONTAINER_NAME")
}

func GetServiceName() string {
	return os.Getenv("SERVICE_NAME")
}

func GetContainerNamespace() string {
	namespaceEnv := os.Getenv("CONTAINER_NAMESPACE")
	if namespaceEnv != "" {
		return namespaceEnv
	}
	return "default"
}

func GetDaemonEndpoint() string {
	return strings.TrimRight(os.Getenv("INTERCEPTOR_DAEMON_HTTP"), "/")
}
