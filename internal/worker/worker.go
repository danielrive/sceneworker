package worker

import (
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// defining a interface for worker

type Worker interface {
	Run(req interface{}) (interface{}, error)
	// define output as a map of strings
}

// Define http worker

type HttpWorker struct {
	HttpClient *http.Client
	Url        string
	HttpMethod string
}

type OutputHttpWorker struct {
	StatusCode  int
	LenghtBody  int64
	ElapsedTime int64
}

func (hw *HttpWorker) Run() (OutputHttpWorker, error) {
	output := OutputHttpWorker{
		StatusCode:  -100,
		LenghtBody:  -100,
		ElapsedTime: -100,
	}
	// Parse url
	urlParsed, err := url.Parse(hw.Url)
	if err != nil {
		log.Fatal("Worker ERROR ", err)
		return output, err

	}

	// Validate data type and procces de body based on the data type

	request, err := http.NewRequest(hw.HttpMethod, urlParsed.String(), nil)
	if err != nil {
		// panic by log
		log.Println("Error creating request", err)
		return output, err
	}
	start := time.Now()
	response, err := hw.HttpClient.Do(request)
	elapsed := time.Since(start)
	// convert to miliseconds

	if err != nil {
		log.Println("Error sending request", err)
		return output, err
	}

	contentLengthStr := response.Header.Get("Content-Length")
	if contentLengthStr != "" {
		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			log.Fatal("Error parsing Content-Length:", err)
			return output, err
		}
		output.LenghtBody = int64(contentLength)
	} else {
		output.LenghtBody = int64(0)
	}

	output.ElapsedTime = elapsed.Nanoseconds() / int64(time.Millisecond)

	output.StatusCode = response.StatusCode

	return output, nil
}
