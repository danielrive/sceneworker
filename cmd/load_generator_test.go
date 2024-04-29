package main

import (
	"devoteam-load-generator/common"
	"devoteam-load-generator/internal/worker"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWorkerGet(t *testing.T) {
	// Define expected response data
	expectedData := worker.OutputHttpWorker{
		StatusCode: 200,
	}
	//expectedData := []byte(`{"message": "success"}`)
	jsonData, err := json.Marshal(expectedData)

	fmt.Println(jsonData)

	if err != nil {
		t.Errorf("Error marshalling mock response: %v", err)
		return
	}
	// Create a Gin router for mocking
	ginRouter := gin.Default()

	// Define a handler for the expected URL and method
	ginRouter.GET("/", func(c *gin.Context) {

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write(jsonData)
		c.Writer.WriteHeader(http.StatusOK)

	})

	ginRouter.POST("/cart", func(c *gin.Context) {

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write(jsonData)
		c.Writer.WriteHeader(http.StatusOK)

	})

	// Create a server using the Gin router
	ts := httptest.NewServer(ginRouter)
	defer ts.Close()

	// define HttpWorker
	client := &http.Client{}
	testGetWorker := worker.HttpWorker{
		Url:        ts.URL,
		HttpMethod: "GET",
		HttpClient: client,
	}

	// Update Worker function to use the test server URL
	_, err = testGetWorker.Run()

	// Assertions
	// compare if error is different to nil or ERROR

	if err != nil && err != common.ErrorTimeout {
		t.Errorf("Unexpected error: %v", err)
	}

}

func TestWorkerPost(t *testing.T) {
	// Define expected response data
	expectedData := worker.OutputHttpWorker{
		StatusCode: 200,
	}
	//expectedData := []byte(`{"message": "success"}`)
	jsonData, err := json.Marshal(expectedData)

	fmt.Println(jsonData)

	if err != nil {
		t.Errorf("Error marshalling mock response: %v", err)
		return
	}
	// Create a Gin router for mocking
	ginRouter := gin.Default()

	ginRouter.POST("/cart", func(c *gin.Context) {

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Write(jsonData)
		c.Writer.WriteHeader(http.StatusOK)

	})

	// Create a server using the Gin router
	ts2 := httptest.NewServer(ginRouter)
	defer ts2.Close()

	// define HttpWorker
	client := &http.Client{}

	testPostWorker := worker.HttpWorker{
		Url:        ts2.URL,
		HttpMethod: "POST",
		HttpClient: client,
	}
	// Update Worker function to use the test server URL
	_, err = testPostWorker.Run()

	// Assertions
	// compare if error is different to nil or ERROR

	if err != nil && err != common.ErrorTimeout {
		t.Errorf("Unexpected error: %v", err)
	}

}
