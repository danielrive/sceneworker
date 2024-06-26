package sceneworker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWorkerGet(t *testing.T) {
	// Define expected response data
	expectedData := OutputHttpWorker{
		StatusCode: 200,
	}
	//expectedData := []byte(`{"message": "success"}`)
	jsonData, err := json.Marshal(expectedData)

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

	// Create a server using the Gin router
	ts := httptest.NewServer(ginRouter)
	defer ts.Close()

	// define HttpWorker
	client := &http.Client{}
	testGetWorker := HttpWorker{
		Url:        ts.URL,
		HttpMethod: "GET",
		HttpClient: client,
	}

	// Update Worker function to use the test server URL
	response, err := testGetWorker.Run()

	// Assertions
	if response.ElapsedTime < 0 {
		t.Errorf("Expected time grather than 0, got %d", response.ElapsedTime)
	}

	// compare if error is different to nil or ERROR

	if err != nil && err != ErrorTimeout {
		t.Errorf("Unexpected error: %v", err)
	}

}

func TestWorkerPost(t *testing.T) {
	// Define expected response data
	expectedData := OutputHttpWorker{
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

	ginRouter.POST("/setCurrency", func(c *gin.Context) {

		c.Writer.Write(jsonData)
		c.Writer.WriteHeader(http.StatusOK)

	})

	// Create a server using the Gin router
	ts2 := httptest.NewServer(ginRouter)
	defer ts2.Close()

	// define HttpWorker
	client := &http.Client{}

	testPostWorker := HttpWorker{
		Url:        ts2.URL,
		HttpMethod: "POST",
		HttpClient: client,
	}
	// Update Worker function to use the test server URL
	_, err = testPostWorker.Run()

	// Assertions
	// compare if error is different to nil or ERROR

	if err != nil && err != ErrorTimeout {
		t.Errorf("Unexpected error: %v", err)
	}

}
