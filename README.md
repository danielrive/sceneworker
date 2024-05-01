# Sceneworker 

A Go package to create workers and Scenarios that can be used by [Boomer](https://github.com/myzhan/boomer) load-generator.Boomer allows the creation of multiple Go routines that will execute the job passed as a function.

## Before to start

### Worker interface

Worker interface defines the contract for worker implementations. Implementations of this interface can be used to execute specific jobs, such as send HTTP request, gRPC(for future), or custom protocols.
Worker implements Run() function that executes a specific job defined by the worker implementation, This function should contain the logic for generating load, sending requests, measuring response times, and generating a respective output.


  ```Go
  type Worker interface {
  	Run(req interface{}) (interface{}, error)
  }
  ```

#### HttpWorker 

HttpWorker is an implementation of the Worker interface for performing HTTP request. 
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
### Scenario interface

A scenario is a group of workers that will be executed in specific order, The idea behind this is be able to define simple workers that can be join together to create complex scenarios.

Scenario interface defines the contract for to Create Scenario. Implementations of this interface can be used to joind specific workers,for instance an HTTP scenario can have the logic to manage cookies or HTTP session.

Scenario implements Run() function that contains the logic to execute the a group of workers and generate outputs that can be used for Boomer.

  ```Go
type Scenario interface {
	Run(req interface{}) (interface{}, error)
}
  ```

#### HttpScenario

HttpScenario is an implementation of the Scenario interface that can be used to run a group of HttpWorkers.

This implements Run() function that get a group of workers and run them in a specific order, passing as the cookies generated for the HTTP request sent.

 ```Go
type HttpScenario struct {
	Name string
	// array of functions
	HttpWorkers []HttpWorker
}

type OutputHttpWorkerScenario struct {
	StatusCode  int   `json:"StatusCode"`
	LenghtBody  int64 `json:"LenghtBody"`
	ElapsedTime int64 `json:"ElapsedTime"`
	Err         error `json:"Err"`
}
```


