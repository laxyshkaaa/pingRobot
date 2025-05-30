package workerPool

import (
	"net/http"
	"time"
)

type Worker struct {
	client *http.Client
}

func NewWorker(time time.Duration) *Worker {
	return &Worker{
		client: &http.Client{
			Timeout: time,
		},
	}
}
func (w Worker) Process(j Job) Result {
	result := Result{URL: j.URL}
	now := time.Now()
	resp, err := w.client.Get(j.URL)
	if err != nil {
		result.Error = err
		return result
	}
	result.StatusCode = resp.StatusCode
	result.ResponseTime = time.Since(now)
	defer resp.Body.Close()
	return result

}
