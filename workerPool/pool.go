package workerPool

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type Job struct {
	URL string
}

type Result struct {
	URL          string
	StatusCode   int
	ResponseTime time.Duration
	Error        error
}

func (r Result) Info() string {
	if r.Error != nil {
		return fmt.Sprintf("Error - [%s], %s", r.URL, r.Error.Error())
	}
	return fmt.Sprintf("URL - [%s], %s, status code: %d", r.URL, r.ResponseTime.String(), r.StatusCode)
}

type Pool struct {
	worker       []*Worker
	workersCount int
	Jobs         chan Job
	time         time.Duration
	Results      chan Result
	wg           *sync.WaitGroup
	stopped      bool
}

func New(workersCount int, time time.Duration, results chan Result) *Pool {
	return &Pool{
		worker:       make([]*Worker, 0, workersCount),
		workersCount: workersCount,
		time:         time,
		Jobs:         make(chan Job, 100),
		Results:      results,
		wg:           new(sync.WaitGroup),
	}
}

func (p *Pool) Init() {
	for i := 0; i < p.workersCount; i++ {
		w := NewWorker(p.time)
		p.worker = append(p.worker, w)
		go p.initWorker(i, w)
	}
}

func (p *Pool) initWorker(id int, w *Worker) {
	for job := range p.Jobs {
		time.Sleep(1 * time.Second)
		p.Results <- w.Process(job)
		p.wg.Done()

	}
	log.Printf("Работник с id %d закончил работу\n", id)
}
func (p *Pool) Push(j Job) {
	if p.stopped {
		return
	}
	p.Jobs <- j
	p.wg.Add(1)
}
func (p *Pool) Stop() {
	p.stopped = true
	close(p.Jobs)
	p.wg.Wait()
}
