package main

import (
	"devoteam-load-generator/common"
	"devoteam-load-generator/internal/worker"
	"devoteam-load-generator/utils"
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

func setCurrency() {

	// create http worker
	serCurrencyWorker := worker.HttpWorker{
		Url:         "http://192.168.0.100:8080/setCurrency",
		HttpMethod:  "POST",
		HttpClient:  client,
		ContentType: "application/x-www-form-urlencoded",
		Body: map[string]string{
			"currency_code": utils.PickupRandom(common.Currencies),
		},
	}

	outputHttpWorker, err := serCurrencyWorker.Run()

	if err != nil {
		globalBoomer.RecordFailure(serCurrencyWorker.Url, err.Error(), outputHttpWorker.ElapsedTime, err.Error())
	} else {
		globalBoomer.RecordSuccess(serCurrencyWorker.Url, strconv.Itoa(outputHttpWorker.StatusCode), outputHttpWorker.ElapsedTime, outputHttpWorker.LenghtBody)
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

	task2 := &boomer.Task{
		Name:   "setCurrency",
		Weight: 5,
		Fn:     setCurrency,
	}

	ts.AddTask(task1)
	ts.AddTask(task2)

	numClients := 2
	spawnRate := float64(10)
	globalBoomer = boomer.NewStandaloneBoomer(numClients, spawnRate)
	globalBoomer.AddOutput(boomer.NewConsoleOutput())
	globalBoomer.Run(task1, task2)
}
