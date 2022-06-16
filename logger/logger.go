package logger

import "log"

func Error(message string, err error) {
	if err != nil {
		log.Fatal("[Error] "+message, err)
	}
}

func Success(message string) {
	log.Println("[Success] " + message)
}
