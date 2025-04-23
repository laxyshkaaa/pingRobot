package main

import (
	"fmt"
	"os"
	"os/signal"
	"pingRobot/workerPool"
	"syscall"
	"time"
)

const (
	INTERVAL        = time.Second * 10
	REQUEST_TIMEOUT = time.Second * 2
	WORKERS_COUNT   = 5
)

var urls = []string{
	"https://workshop.zhashkevych.com/",
	"https://golang-ninja.com/",
	"https://ru.wikipedia.org/wiki/%D0%97%D0%B0%D0%B3%D0%BB%D0%B0%D0%B2%D0%BD%D0%B0%D1%8F_%D1%81%D1%82%D1%80%D0%B0%D0%BD%D0%B8%D1%86%D0%B0",
	"https://google.com/",
	"https://golang.org/",
}

func main() {
	results := make(chan workerPool.Result)
	WorkerPools := workerPool.New(WORKERS_COUNT, REQUEST_TIMEOUT, results)
	WorkerPools.Init()
	go generateJobs(WorkerPools)
	go processResults(results)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	WorkerPools.Stop()
}
func processResults(results chan workerPool.Result) {
	go func() {
		for result := range results {
			fmt.Println("result:", result.Info())
		}
	}()
}
func generateJobs(wp *workerPool.Pool) {
	for {
		for _, url := range urls {
			wp.Push(workerPool.Job{URL: url})
		}
		time.Sleep(INTERVAL)
	}
}
