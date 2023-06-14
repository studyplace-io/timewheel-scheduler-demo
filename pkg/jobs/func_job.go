package jobs

import (
	"context"
	"fmt"
)

// Function represents an argument-less function which returns a generic type R and a possible error.
type Function func(context.Context) (interface{}, error)

// FunctionJob represents a Job that invokes the passed Function, implements the quartz.Job interface.
type FunctionJob struct {
	JobStatus     JobStatus
	Type          JobType
	Times         int
	function      Function
	desc          string
	LatestResult  *JobCommandResult
	HistoryResult chan *JobCommandResult
	Error         error
}

type FuncResult struct {
	Result interface{}
	Error  error
}

// NewFunctionJob returns a new FunctionJob without an explicit description.
func NewFunctionJob(function Function, times int) JobCommand {
	return &FunctionJob{
		function:     function,
		desc:         fmt.Sprintf("FunctionJob:%p", &function),
		LatestResult: &JobCommandResult{},
		Error:        nil,
		JobStatus:    NotStarted,
		Times:        times,
		Type:         FuncType,
	}
}

// NewFunctionJobWithDesc returns a new FunctionJob with an explicit description.
func NewFunctionJobWithDesc(desc string, function Function) JobCommand {
	return &FunctionJob{
		function:     function,
		desc:         desc,
		LatestResult: &JobCommandResult{},
		Error:        nil,
		JobStatus:    NotStarted,
		Type:         FuncType,
	}
}

// Description returns the description of the FunctionJob.
func (f *FunctionJob) Description() string {
	return f.desc
}

// ShowJobType returns the description of the CurlJob.
func (f *FunctionJob) ShowJobType() string {
	return fmt.Sprintf("Job Type: %s", f.Type)
}

func (f *FunctionJob) SetJobExecuteTime() int {
	return f.Times
}

// Execute is called by a Scheduler when the Trigger associated with this jobs fires.
// It invokes the held function, setting the results in Result and Error members.
func (f *FunctionJob) Execute(ctx context.Context) {

	if f.Times != -1 && f.JobStatus == NotStarted {
		f.HistoryResult = make(chan *JobCommandResult, f.SetJobExecuteTime())
	}

	f.JobStatus = Running
	result, err := f.function(ctx)

	if err != nil {
		f.JobStatus = Fail
		f.LatestResult.Result = nil
		f.LatestResult.Error = err
	} else {
		f.JobStatus = Success
		f.LatestResult.Result = &result
		f.LatestResult.Error = nil
	}

	if f.Times != -1 {
		f.HistoryResult <- f.LatestResult
	}

}

func (f *FunctionJob) GetLatestResult() *JobCommandResult {
	return f.LatestResult
}

func (f *FunctionJob) GetHistoryResult() <-chan *JobCommandResult {
	return f.HistoryResult
}
