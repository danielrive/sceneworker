package main

import (
	"devoteam-load-generator/common"
	"devoteam-load-generator/internal/worker"
	"devoteam-load-generator/utils"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/myzhan/boomer"
)

var client *http.Client
var timeout int = 3
var url string

func index() {

	// create http worker
	indexWorker := worker.HttpWorker{
		Url:        url + "/",
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
		Url:         url + "/setCurrency",
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

func browseProduct() {

	// create http worker
	setCurrencyWorker := worker.HttpWorker{
		Url:         url + "/product/" + utils.PickupRandom(common.Products),
		HttpMethod:  "GET",
		HttpClient:  client,
		ContentType: "application/x-www-form-urlencoded",
	}

	outputHttpWorker, err := setCurrencyWorker.Run()

	if err != nil {
		globalBoomer.RecordFailure(setCurrencyWorker.Url, err.Error(), outputHttpWorker.ElapsedTime, err.Error())
	} else {
		globalBoomer.RecordSuccess(setCurrencyWorker.Url, strconv.Itoa(outputHttpWorker.StatusCode), outputHttpWorker.ElapsedTime, outputHttpWorker.LenghtBody)
	}
}

func checkout() {

	// create http worker
	checkoutWorker := worker.HttpWorker{
		Url:         url + "/cart/checkout",
		HttpMethod:  "POST",
		HttpClient:  client,
		ContentType: "application/x-www-form-urlencoded",
		Body:        utils.FakeCheckout(),
	}

	outputHttpWorker, err := checkoutWorker.Run()

	if err != nil {
		globalBoomer.RecordFailure(checkoutWorker.Url, err.Error(), outputHttpWorker.ElapsedTime, err.Error())
	} else {
		globalBoomer.RecordSuccess(checkoutWorker.Url, strconv.Itoa(outputHttpWorker.StatusCode), outputHttpWorker.ElapsedTime, outputHttpWorker.LenghtBody)
	}
}

var globalBoomer *boomer.Boomer

func main() {
	url = os.Getenv("URL")
	if url == "" {
		panic("URL cannot be empty")
	}
	client = &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	// Scenario
	// 1. index
	// 2. setCurrency
	// 3. browseProducts
	// 4. checkout

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

	task3 := &boomer.Task{
		Name:   "browseProducts",
		Weight: 10,
		Fn:     browseProduct,
	}
	task4 := &boomer.Task{
		Name:   "checkout",
		Weight: 10,
		Fn:     checkout,
	}

	numClients := 100
	spawnRate := float64(200)
	globalBoomer = boomer.NewStandaloneBoomer(numClients, spawnRate)
	//globalBoomer.AddOutput(boomer.NewConsoleOutput())

	globalBoomer.Run(task1, task2, task3, task4)

}
