package main

import (
	"devoteam_load_generator/internal/worker"
	"net/http"
	"strconv"
	"time"

	"github.com/myzhan/boomer"
)

var client *http.Client
var timeout int = 1

func index() {

	// create http worker
	indexWorker := worker.HttpWorker{
		Url:        "http://192.168.0.100:8080/",
		HttpMethod: "GET",
		HttpClient: client,
	}

	outputHttpWorker, err := indexWorker.Run()

	if err != nil {
		globalBoomer.RecordFailure(indexWorker.Url, err.Error(), outputHttpWorker.ElapsedTime, err.Error())
	} else {
		globalBoomer.RecordSuccess(indexWorker.Url, strconv.Itoa(outputHttpWorker.StatusCode), outputHttpWorker.ElapsedTime, outputHttpWorker.LenghtBody)
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
