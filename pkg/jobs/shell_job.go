package jobs

import (
	"context"
	"fmt"
	"os/exec"
)

// ShellJob 执行shell命令Job
type ShellJob struct {
	JobStatus     JobStatus
	Type          JobType
	Times         int
	Cmd           string
	LatestResult  *JobCommandResult
	HistoryResult chan *JobCommandResult
}

// NewShellJob returns a new ShellJob
func NewShellJob(cmd string, times int) JobCommand {
	return &ShellJob{
		Cmd:          cmd,
		LatestResult: &JobCommandResult{},
		Times:        times,
		JobStatus:    NotStarted,
		Type:         ShellType,
	}
}

// Description returns the description of the ShellJob.
func (sh *ShellJob) Description() string {
	return fmt.Sprintf("ShellJob: %s", sh.Cmd)
}

// ShowJobType returns the description of the ShellJob.
func (sh *ShellJob) ShowJobType() string {
	return fmt.Sprintf("Job Type: %s", sh.Type)
}

// Execute 执行Job
func (sh *ShellJob) Execute(ctx context.Context) {

	// 如果是限制次数，需要初始化chan，
	// 实现历史记录存在chan中，用户可以查询
	if sh.Times != -1 && sh.JobStatus == NotStarted {
		sh.HistoryResult = make(chan *JobCommandResult, sh.SetJobExecuteTime())
	}

	// 执行命令
	out, err := exec.CommandContext(ctx, "sh", "-c", sh.Cmd).Output()
	// 修改状态
	sh.JobStatus = Running
	if err != nil {
		sh.JobStatus = Fail
		sh.LatestResult.Error = err
		return
	}

	sh.JobStatus = Success
	sh.LatestResult.Result = string(out)

	// 存到chan中
	if sh.Times != -1 {
		sh.HistoryResult <- sh.LatestResult
	}
}

func (sh *ShellJob) SetJobExecuteTime() int {
	return sh.Times
}

func (sh *ShellJob) GetLatestResult() *JobCommandResult {
	return sh.LatestResult
}

func (sh *ShellJob) GetHistoryResult() <-chan *JobCommandResult {
	return sh.HistoryResult
}
