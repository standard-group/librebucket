package worker

import (
	"log"
)

// Worker processes jobs from the Jobs channel.
type Worker struct {
	ID   int
	Jobs <-chan Job
	Quit chan bool
}

func NewWorker(id int, jobs <-chan Job) *Worker {
	return &Worker{
		ID:   id,
		Jobs: jobs,
		Quit: make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case job := <-w.Jobs:
				if err := job.Run(); err != nil {
					log.Printf("Worker %d: job failed: %v", w.ID, err)
				}
			case <-w.Quit:
				log.Printf("Worker %d stopping", w.ID)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.Quit <- true
}
