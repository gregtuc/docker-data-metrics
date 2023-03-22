package database

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Log struct {
	ContainerIP string  `json:"containerIp"`
	Timestamp   int64   `json:"timestamp"`
	CPUPercent  float64 `json:"cpuPercent"`
	RamPercent  float64 `json:"ramPercent"`
}

var (
	mu sync.Mutex
)

func Write(log Log) {
	mu.Lock()
	existingData := Read()
	existingData = append(existingData, log)

	result, err := json.Marshal(existingData)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("./database/logs.json", result, 0644)
	if err != nil {
		panic(err)
	}
	mu.Unlock()
	return
}

func Read() []Log {
	var data []Log
	file, err := ioutil.ReadFile("./database/logs.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		panic(err)
	}

	return data
}
