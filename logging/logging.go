package logging

import (
	"log"
	"os"
	"strconv"
	"time"
)

func LogToFile(msg string, level string) {
	logEnv := os.Getenv("LOGGING_PATH")
	if logEnv == "" {
		panic("Couldn't find LOG_PATH env")
	}
	file, err := os.OpenFile(logEnv, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	log.SetOutput(file)
	if level == "fatal" {
		log.Fatal(msg)
	} else {
		log.Println(msg)
	}
}

func LogToSnapshotFile(msg string) {
	crLogging := os.Getenv("CR_LOGGING_PATH")
	if crLogging == "" {
		panic("Couldn't find CR_LOGGING_PATH env")
	}
	currentTimestamp := time.Now().UnixNano()
	f, err := os.OpenFile(crLogging,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	stringToStore := msg + "@" + strconv.FormatInt(currentTimestamp, 10)
	if _, err := f.WriteString(stringToStore); err != nil {
		log.Println(err)
	}
}
