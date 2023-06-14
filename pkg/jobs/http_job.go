package jobs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
)

// CurlJob represents a cURL command Job, implements the quartz.Job interface.
type CurlJob struct {
	RequestMethod string
	URL           string
	Body          string
	Headers       map[string]string
	LatestResult  *JobCommandResult
	HistoryResult chan *JobCommandResult
	Type          JobType
	JobStatus     JobStatus
	Times         int
	request       *http.Request
}

type HttpResponse struct {
	Response   string
	StatusCode int
}

// NewCurlJob returns a new CurlJob.
func NewCurlJob(method string, url string, body string, headers map[string]string, times int) (JobCommand, error) {
	_body := bytes.NewBuffer([]byte(body))
	req, err := http.NewRequest(method, url, _body)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	return &CurlJob{
		RequestMethod: method,
		URL:           url,
		Body:          body,
		Headers:       headers,
		LatestResult:  &JobCommandResult{},
		JobStatus:     NotStarted,
		Times:         times,
		request:       req,
		Type:          CullType,
	}, nil
}

// Description returns the description of the CurlJob.
func (cu *CurlJob) Description() string {
	return fmt.Sprintf("CurlJob: %s %s %s", cu.RequestMethod, cu.URL, cu.Body)
}

// ShowJobType returns the description of the CurlJob.
func (cu *CurlJob) ShowJobType() string {
	return fmt.Sprintf("Job Type: %s", cu.Type)
}

func (cu *CurlJob) SetJobExecuteTime() int {
	return cu.Times
}

// Execute is called by a Scheduler when the Trigger associated with this jobs fires.
func (cu *CurlJob) Execute(ctx context.Context) {
	if cu.Times != -1 && cu.JobStatus == NotStarted {
		cu.HistoryResult = make(chan *JobCommandResult, cu.SetJobExecuteTime())
	}

	client := &http.Client{}
	cu.request = cu.request.WithContext(ctx)
	cu.JobStatus = Running
	resp, err := client.Do(cu.request)
	if err != nil {
		cu.JobStatus = Fail
		cu.LatestResult.Result = nil
		cu.LatestResult.Error = err
		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		cu.JobStatus = Success
	} else {
		cu.JobStatus = Fail
	}

	cu.LatestResult.Result = string(body)
	cu.LatestResult.Error = nil
	if cu.Times != -1 {
		cu.HistoryResult <- cu.LatestResult
	}
}

func (cu *CurlJob) GetLatestResult() *JobCommandResult {
	return cu.LatestResult
}

func (cu *CurlJob) GetHistoryResult() <-chan *JobCommandResult {
	return cu.HistoryResult
}
