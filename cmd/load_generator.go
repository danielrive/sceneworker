package main

import (
	"devoteam-load-generator/common"
	"devoteam-load-generator/internal/scenario"
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

func scenario1() {

	indexWorker := worker.HttpWorker{
		Url:        url + "/",
		HttpMethod: "GET",
		HttpClient: client,
	}

	serCurrencyWorker := worker.HttpWorker{
		Url:         url + "/setCurrency",
		HttpMethod:  "POST",
		HttpClient:  client,
		ContentType: "application/x-www-form-urlencoded",
		Body: map[string]string{
			"currency_code": utils.PickupRandom(common.Currencies),
		},
	}

	browseProductWorker := worker.HttpWorker{
		Url:         url + "/product/" + utils.PickupRandom(common.Products),
		HttpMethod:  "GET",
		HttpClient:  client,
		ContentType: "application/x-www-form-urlencoded",
	}

	viewCartWorker := worker.HttpWorker{
		Url:         url + "/cart",
		HttpMethod:  "GET",
		HttpClient:  client,
		ContentType: "application/x-www-form-urlencoded",
	}

	addToCartWorker := worker.HttpWorker{
		Url:         url + "/cart",
		HttpMethod:  "POST",
		HttpClient:  client,
		ContentType: "application/x-www-form-urlencoded",
		Body: map[string]string{
			"product_id": utils.PickupRandom(common.Products),
			"quantity":   "1",
		},
	}
	checkoutWorker := worker.HttpWorker{
		Url:         url + "/cart/checkout",
		HttpMethod:  "POST",
		HttpClient:  client,
		ContentType: "application/x-www-form-urlencoded",
		Body:        utils.FakeCheckout(),
	}

	scenario1 := scenario.HttpScenario{
		Name: "scenario1",
		HttpWorkers: []worker.HttpWorker{
			indexWorker,
			serCurrencyWorker,
			browseProductWorker,
			viewCartWorker,
			addToCartWorker,
			checkoutWorker,
		},
	}

	outputScearnio := scenario1.Run()

	if outputScearnio.Err != nil {
		globalBoomer.RecordFailure(scenario1.Name, outputScearnio.Err.Error(), outputScearnio.ElapsedTime, outputScearnio.Err.Error())
	} else {
		globalBoomer.RecordSuccess(scenario1.Name, strconv.Itoa(outputScearnio.StatusCode), outputScearnio.ElapsedTime, outputScearnio.LenghtBody)
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

	task1 := &boomer.Task{
		Name:   "index",
		Weight: 15,
		Fn:     scenario1,
	}

	numClients := 1
	spawnRate := float64(0.5)
	globalBoomer = boomer.NewStandaloneBoomer(numClients, spawnRate)

	globalBoomer.Run(task1)

}
