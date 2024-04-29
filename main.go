package main

import (
	"devoteam-load-generator/internal/worker"
	"net/http"
	"strconv"
	"time"

	"github.com/myzhan/boomer"
)

var client *http.Client
var timeout int = 1

func index() {

	// create http worker
	index_worker := worker.HttpWorker{
		Url:        "http://192.168.0.100:8080/",
		HttpMethod: "GET",
		HttpClient: client,
	}

	outputHttpWorker, err := index_worker.Run()

	if err != nil {
		globalBoomer.RecordFailure(index_worker.Url, err.Error(), outputHttpWorker.ElapsedTime, err.Error())
	} else {
		globalBoomer.RecordSuccess(index_worker.Url, strconv.Itoa(outputHttpWorker.StatusCode), outputHttpWorker.ElapsedTime, outputHttpWorker.LenghtBody)
	}
}

var globalBoomer *boomer.Boomer

func main() {
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
