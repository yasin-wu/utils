package dcron

import "context"

// TaskExecutor 任务执行器接口
type TaskExecutor interface {
	Name() string
	Execute(ctx context.Context, task *Task) error
}
