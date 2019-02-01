package pkg

import (
	//"fmt"
	//"os"
	//"strconv"
)

var (
	MaxWorker = 10//os.Getenv("MAX_WORKERS")
	MaxQueue  = 10//os.Getenv("MAX_QUEUE")
)

type JobWithResultChan struct {
	Job
	Result chan interface{}
}

type Dispatcher struct {
	MaxWorkers int

	// A buffered channel that we can send work requests on.
	JobQueue chan Job

	// A pool of workers channels that are registered with the dispatcher
	WorkerPool chan chan Job
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		MaxWorkers: int(MaxWorker),
		JobQueue: make(chan Job, MaxQueue),
		WorkerPool: make(chan chan Job, MaxWorker),
	}
}

func (d *Dispatcher) Submit(job Job) {
	d.JobQueue <- job
}

func (d *Dispatcher) Run() {
	// starting n number of workers
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <- d.JobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.WorkerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
