package main

import (
	"encoding/json"
	"log"
	"os"
)

type jsonAccess struct {
	filename string
}

func (j *jsonAccess) openJSONFile() *os.File {
	unreadJSONFile, err := os.OpenFile(j.filename, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalf("Unable to open JSON file: %v", err)
	}
	return unreadJSONFile
}

func (j *jsonAccess) readConfig(conf *config) int64 {
	file := j.openJSONFile()
	defer file.Close()

	json.NewDecoder(file).Decode(conf)

	return j.getTimestamp()
}

func (j *jsonAccess) readIDs(ids *[]string) int64 {
	file := j.openJSONFile()
	defer file.Close()

	json.NewDecoder(file).Decode(ids)

	return j.getTimestamp()
}

func (j *jsonAccess) write(data interface{}) int64 {
	messageJSON, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Unable to encode JSON: %v", err)
	}

	file := j.openJSONFile()
	defer file.Close()

	file.Truncate(0)
	file.Write(messageJSON)
	if err != nil {
		log.Fatalf("Unable to write JSON file: %v", err)
	}

	return j.getTimestamp()
}

func (j *jsonAccess) getTimestamp() int64 {
	info, err := os.Stat(j.filename)
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		log.Fatalf("Unable to get file status: %v", err)
	}
	return info.ModTime().Unix()
}
