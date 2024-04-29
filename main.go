package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/myzhan/boomer"
)

var client *http.Client
var timeout int = 1

func index() {
	request, err := http.NewRequest("GET", "http://192.168.0.100:8080", nil)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	startTime := time.Now()
	response, err := client.Do(request)
	elapsed := time.Since(startTime)
	if err != nil {
		globalBoomer.RecordFailure("Get /", "error", elapsed.Nanoseconds()/int64(time.Millisecond), err.Error())
	} else {
		globalBoomer.RecordSuccess("Get /", strconv.Itoa(response.StatusCode),
			elapsed.Nanoseconds()/int64(time.Millisecond), response.ContentLength)

		response.Body.Close()
	}
}

var globalBoomer *boomer.Boomer

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client = &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	ts := boomer.NewWeighingTaskSet()

	task1 := &boomer.Task{
		Name:   "index",
		Weight: 20,
		Fn:     index,
	}

	ts.AddTask(task1)

	numClients := 2
	spawnRate := float64(10)
	globalBoomer = boomer.NewStandaloneBoomer(numClients, spawnRate)
	globalBoomer.AddOutput(boomer.NewConsoleOutput())
	globalBoomer.Run(task1)
}
