package model

import (
	"log"
	"os"
	"time"
)

type Job struct {
	Data []byte
}

func (j Job) Process() error {
	dataString := string(j.Data)

	dataWithTimestamp := dataString + " - Processed at " + time.Now().Format(time.RFC3339)

	log.Printf("Processed data: %s\n", dataWithTimestamp)

	if err := writeToFile(dataString); err != nil {
		log.Printf("failed to write to file: %s", err.Error())
	}

	return nil
}

func writeToFile(data string) error {
	file, err := os.OpenFile("processed_data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data + "\n")
	if err != nil {
		return err
	}

	return nil
}
