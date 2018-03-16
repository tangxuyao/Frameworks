package Frameworks

import "github.com/prometheus/common/log"

var (
	MaxWorker = 10
	MaxQueue  = 10
)

type Payload interface {
	Do() error
}

type Job struct {
	Payload Payload
}

var JobQueue chan Job

type Worker struct {
	WorkerPool chan chan Job
	JobChannel chan Job
	quit       chan bool
}

func NewWorker(pool chan chan Job) Worker {
	return Worker{
		WorkerPool: pool,
		JobChannel: make(chan Job),
		quit:       make(chan bool)}
}

func (w Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				if err := job.Payload.Do(); err != nil {
					log.Errorf("Error Job Do:%s", err.Error())
				}
			case <-w.quit:
				return
			}
		}
	}()
}


func (w Worker) Stop() {
	go func() {
		w.quit <- true
	}()
}