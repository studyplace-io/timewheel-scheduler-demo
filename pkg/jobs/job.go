package jobs

import "context"

// JobCommand Job接口对象
type JobCommand interface {
	// Execute 执行命令
	Execute(ctx context.Context)
	// Description 描述
	Description() string
	// ShowJobType Job类型
	ShowJobType() string
	// SetJobExecuteTime 给timewheel设置执行次数，用户不需要用
	SetJobExecuteTime() int

	// 使用方调用的方法
	// GetLatestResult 获取Job最新的执行结果
	GetLatestResult() *JobCommandResult
	// GetHistoryResult 存储Job执行的结果，只有限定次数的Job才会存储
	GetHistoryResult() <-chan *JobCommandResult
}

type JobCommandResult struct {
	Result interface{}
	Error  error
}

// JobStatus Job执行状态
type JobStatus string

const (
	NotStarted = "not started"
	Running    = "running"
	Success    = "success"
	Fail       = "fail"
)

// JobType Job种类
type JobType string

const (
	ShellType = "shell"
	CullType  = "curl"
	FuncType  = "func"
)
