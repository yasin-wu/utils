package dcron

// Task 任务定义
type Task struct {
	ID          string         `json:"id"`
	CronJobID   int            `json:"cron_job_id"`
	Name        string         `json:"name"`
	Expression  string         `json:"expression"` // cron表达式
	Executor    string         `json:"executor"`   // 执行器名称
	Parameters  map[string]any `json:"parameters"`
	Description string         `json:"description"`
	Enabled     bool           `json:"enabled"`
}
