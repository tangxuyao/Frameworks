package Frameworks

type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
}

func NewDispatcher(max int) *Dispatcher {
	JobQueue = make(chan Job)
	pool := make(chan chan Job, max)
	return &Dispatcher{WorkerPool: pool, MaxWorkers: max}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-JobQueue:

			go func(job Job) {
				channel := <-d.WorkerPool
				channel <- job
			}(job)
		}
	}
}
