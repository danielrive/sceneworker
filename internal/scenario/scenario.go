package scenario

import (
	"devoteam-load-generator/internal/worker"
)

type HttpScenario struct {
	Name string
	// array of functions
	HttpWorkers []worker.HttpWorker
}

type OutputHttpWorkerScenario struct {
	StatusCode  int   `json:"StatusCode"`
	LenghtBody  int64 `json:"LenghtBody"`
	ElapsedTime int64 `json:"ElapsedTime"`
	Err         error `json:"Err"`
}

type Scenario2 interface {
	Run(req interface{}) (interface{}, error)
}

func (s *HttpScenario) Run() OutputHttpWorkerScenario {

	OutputScenario := OutputHttpWorkerScenario{
		StatusCode:  -100,
		LenghtBody:  -100,
		ElapsedTime: -100,
		Err:         nil,
	}

	onStart, _ := s.HttpWorkers[0].Run()
	for w := 1; w < len(s.HttpWorkers); w++ {

		s.HttpWorkers[w].Cookies = onStart.Cookies
		response, err := s.HttpWorkers[w].Run()
		OutputScenario.StatusCode = response.StatusCode
		OutputScenario.LenghtBody = response.LenghtBody
		OutputScenario.ElapsedTime = response.ElapsedTime
		OutputScenario.Err = err

	}
	return OutputScenario
}
