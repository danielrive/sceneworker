package worker

import (
	"bytes"
	"devoteam-load-generator/common"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// defining a interface for worker

type Worker interface {
	Run(req interface{}) (interface{}, error)
	// define output as a map of strings
}

// Define http worker

type HttpWorker struct {
	HttpClient  *http.Client
	Url         string
	HttpMethod  string
	ContentType string
	Body        map[string]string
}

type OutputHttpWorker struct {
	StatusCode  int   `json:"StatusCode"`
	LenghtBody  int64 `json:"LenghtBody"`
	ElapsedTime int64 `json:"ElapsedTime"`
}

func (hw *HttpWorker) Run() (OutputHttpWorker, error) {
	// empty data by default
	data := []byte{}

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

	// Validate if Method is POST to setup the datatype of the answer
	if hw.HttpMethod == "POST" {
		if hw.ContentType == "application/json" {
			data, _ = json.Marshal(hw.Body)
		}
		if hw.ContentType == "application/x-www-form-urlencoded" {

			htmlForm := url.Values{}
			for key, value := range hw.Body {
				htmlForm.Set(key, value)
			}
			data = []byte(htmlForm.Encode())
		}

	}

	request, err := http.NewRequest(hw.HttpMethod, urlParsed.String(), bytes.NewBuffer(data))
	request.Header.Set("Content-Type", hw.ContentType)
	// Validate if error is Timeout or the site is not up
	if err != nil {
		log.Fatal("Worker ERROR, creating request:", err)
	}
	start := time.Now()
	response, err := hw.HttpClient.Do(request)
	elapsed := time.Since(start)

	// Validate if error is Timeout or the site is not up
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout") {
			//log.Println("Error creating request", err)
			return output, common.ErrorTimeout
		} else {
			log.Fatal("Worker ERROR, bad request ", err)
			return output, err
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
