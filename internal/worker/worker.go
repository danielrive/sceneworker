package worker

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Common errors
var (
	errorTimeout = errors.New("Timeout")
	errorRefused = errors.New("Refused")
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

	// Validate if error is Timeout or the site is not up
	if err != nil {
		log.Fatal("Worker ERROR, creating request:", err)
	}
	start := time.Now()
	response, err := hw.HttpClient.Do(request)
	elapsed := time.Since(start)

	// Validate if error is Timeout or the site is not up
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			//log.Println("Error creating request", err)
			return output, errorTimeout
		} else {
			//log.Println("Error creating request", err)
			return output, errorRefused
		}
	}
	// convert to miliseconds
	contentLengthStr := response.Header.Get("Content-Length")
	// validate content length
	if contentLengthStr != "" {
		contentLength, err := strconv.Atoi(contentLengthStr)
		if err != nil {
			log.Fatal("Worker ERROR, Error parsing Content-Length:", err)
		}
		output.LenghtBody = int64(contentLength)
	} else {
		output.LenghtBody = int64(0)
	}

	output.ElapsedTime = elapsed.Nanoseconds() / int64(time.Millisecond)

	output.StatusCode = response.StatusCode

	return output, nil
}
