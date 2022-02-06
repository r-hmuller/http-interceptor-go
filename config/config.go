package config

import (
	"httpInterceptor/logging"
	"os"
	"strconv"
)

func GetApplicationURL() string {
	return os.Getenv("HOST")
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
