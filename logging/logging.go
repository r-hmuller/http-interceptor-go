package logging

import (
	"log"
	"os"
)

func LogToFile(msg string, level string) {
	logEnv := os.Getenv("LOG_PATH")
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
