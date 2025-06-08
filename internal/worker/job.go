package worker

// Job represents a unit of work to be processed by a worker.
type Job interface {
	Run() error
}

// ExampleJob is a sample implementation of the Job interface.
type ExampleJob struct {
	Payload string
}

func (e *ExampleJob) Run() error {
	// Implement job logic here
	return nil
}