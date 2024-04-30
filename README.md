# Load Generator for microservices-demo

# Overview

This project is a load generator written in Go using the [Boomer package](https://github.com/myzhan/boomer). Boomer allows the creation of multiple Go routines that will execute the job passed as a function.

In this version of the load generator, I've integrated improvements to enhance the way HTTP requests are handled, providing more flexibility and reliability in load-testing scenarios.

The scenario created for this project is to generate load in the microservices defined in [microservices-demo project](https://github.com/GoogleCloudPlatform/microservices-demo) but can be modified to run against another HTTP endpoints.


## Before to start

### Worker interface

Worker interface defines the contract for load test worker implementations. Implementations of this interface can be used to create specific load test scenario,such as HTTP, gRPC, or custom protocols.
Worker implements Run() function that executes load test scenario defined by the worker implementation, This function should contain the logic for generating load, sending requests, measuring response times, and generating and respective output.

  ```Go
  type Worker interface {
  	Run(req interface{}) (interface{}, error)
  }
  ```
### HttpWorker 

HttpWorker is an implementation of the Worker interface for performing HTTP load tests. 
This implements Run() function creating and sending the HTTP request to the endpoint specified.

```Go
type HttpWorker struct {
    // HttpClient is the HTTP client used for sending requests.
    HttpClient *http.Client
    // Url is the URL to which requests will be sent.
    Url string
    // HttpMethod is the HTTP method (e.g., GET, POST) used for the requests.
    HttpMethod string
    // ContentType specifies the content type of the request body.
    ContentType string
    // Body is a map containing key-value pairs for the request body parameters.
    Body map[string]string
}
```

### Defining Boomer tasks

Boomer task represents a specific operation or action that you want to perform repeatedly under load, for example in our case execute Run() function for specific HttpWorker

Bommer defines the task as this
```Go
type Task struct {
	// The weight of a task determines the relative proportion of goroutines that should be assigned to execute that task
  // compared to other tasks. Higher-weight tasks receive more goroutines, while lower-weight tasks receive fewer goroutines
	Weight int
	// Fn is called by the goroutines allocated to this task, in a loop.
	Fn   func()
	Name string
}
```
### Creating Task

Before creating a Boomer task you must define the function that the task will execute. An HttpWorker can be implemented, this abstracts the job of sending HTTP request and generating an output that Boomer will use to show results.

#### Creating a function that implements a HttpWorker.

setCurrency function creates an HttpWorker with specific attributes like, Http method and body to send in the Post Request. 
Run() function will return information like Http code in the response, the time elapsed in the request, and content-lengt.This information is passed to Boomer to record a successful/failed test.
```Go
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
```

# Usage


## Run via Docker

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/yourusername/load-generator.git
   cd load-generator

## Binary
To use this load generator, make sure you have Go installed on your system. You can download and install it from the official Go website.
