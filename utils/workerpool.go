package utils

import (
	"fmt"
	"reflect"
)

type IWorkerPool interface {
	// Start starts consuming jobs from the queue
	Start()
	// Stop tears down the worker pool
	Stop()
	// Enqueue adds a job to wait queue
	Enqueue(job func(...any) error, args ...any) error
}

type IJob interface {
	// Run execute the job and return error or nil
	run() error
	// TODO: support onError executions
}

type Job struct {
	Job         any
	Args        []any
	PostJob     any   // optional
	PostJobArgs []any // optional
}

func NewJob(mainJob any, args []any, postJob any, postJobArgs []any) (*Job, error) {

	if mainJob == nil {
		return nil, fmt.Errorf("main job cannot be nil")
	}

	if err := validateFunction(mainJob, args...); err != nil {
		return nil, fmt.Errorf("invalid main job: %w", err)
	}

	return &Job{
		Job:         mainJob,
		Args:        args,
		PostJob:     postJob,
		PostJobArgs: postJobArgs,
	}, nil
}

func (j *Job) run() error {
	jobReturns := execFunction(j.Job, assertInterfacesToValues(j.Args))

	if j.PostJob != nil {
		// Append main job's output to post job args
		allArgs := append(assertInterfacesToValues(j.PostJobArgs), jobReturns...)

		execFunction(j.PostJob, allArgs)
	}
	return nil
}

// execFunction is a helper to execute any function with its arguments using reflection
func execFunction(fn any, args []reflect.Value) []reflect.Value {
	fnValue := reflect.ValueOf(fn)

	// Call the function
	results := fnValue.Call(args)

	return results
}

func assertInterfacesToValues(values []any) []reflect.Value {
	var out = make([]reflect.Value, len(values))

	for i, val := range values {
		out[i] = reflect.ValueOf(val)
	}
	return out
}

type WorkerPool struct {
	// the number of workers
	Size       uint
	WaitQueue  chan Job
	ExitSignal chan struct{}
}

func NewWorkerPool(size uint, queueSize uint) WorkerPool {
	return WorkerPool{
		Size:       size,
		WaitQueue:  make(chan Job, queueSize),
		ExitSignal: make(chan struct{}),
	}
}

func (p *WorkerPool) Start() {
	for i := uint(0); i < p.Size; i++ {
		go func(workerID uint) {
			for {
				select {
				// if the exit signal is initiated
				case <-p.ExitSignal:
					// end the worker go routine
					fmt.Printf("Exiting worker with ID: %d\n", workerID)
					return
				// if there is a new job on the queue
				case job, ok := <-p.WaitQueue:
					if !ok {
						fmt.Printf("Exiting worker %d because the wait channel is closed", workerID)
						return
					}
					if err := job.run(); err != nil {
						fmt.Printf("error while executing job %v", err)
					}
				}
			}
		}(i)
	}

}

func (p *WorkerPool) Stop() {
	// send a signal to all go routines to exit
	for i := uint(0); i < p.Size; i++ {
		p.ExitSignal <- struct{}{}
	}
}

func (p *WorkerPool) Enqueue(job *Job) error {
	// TODO: implement timeout when enqueueing
	p.WaitQueue <- *job
	return nil
}

func validateFunction(fn any, args ...any) error {
	funcType := reflect.TypeOf(fn)

	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("job must be a function")
	}

	if len(args) != funcType.NumIn() {
		return fmt.Errorf("number of arguments mismatch: expected %d, got %d", funcType.NumIn(), len(args))
	}

	// Type checking
	for i, arg := range args {
		argType := reflect.TypeOf(arg)
		expectedArgType := funcType.In(i)

		if !argType.AssignableTo(expectedArgType) {
			return fmt.Errorf("argument %d type mismatch: got %v, expected %v", i, argType, expectedArgType)
		}
	}

	return nil
}
