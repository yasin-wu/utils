package dcron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// DistributedCron 分布式定时任务管理器
type DistributedCron struct {
	cron         *cron.Cron
	locker       DistributedLocker
	tasks        map[string]*Task
	executors    map[string]TaskExecutor
	instanceID   string
	mutex        sync.RWMutex
	config       *Config
	runningTasks sync.Map
	logger       Logger
}

// Config 配置
type Config struct {
	InstanceID          string        `json:"instance_id"`
	LockTTL             time.Duration `json:"lock_ttl"`
	LockRefreshInterval time.Duration `json:"lock_refresh_interval"`
	EnableDistributed   bool          `json:"enable_distributed"`
}

// NewDistributedCron 创建分布式定时任务管理器
func NewDistributedCron(config *Config, locker DistributedLocker, logger Logger) *DistributedCron {
	if config == nil {
		config = &Config{
			InstanceID:          fmt.Sprintf("instance-%d", time.Now().UnixNano()),
			LockTTL:             30 * time.Second,
			LockRefreshInterval: 10 * time.Second,
			EnableDistributed:   true,
		}
	}
	if logger == nil {
		logger = NewDefaultLogger()
	}
	secondParser := cron.NewParser(cron.Second | cron.Minute |
		cron.Hour | cron.Dom | cron.Month | cron.DowOptional | cron.Descriptor)
	c := &DistributedCron{
		cron:       cron.New(cron.WithParser(secondParser), cron.WithChain()),
		locker:     locker,
		tasks:      make(map[string]*Task),
		executors:  make(map[string]TaskExecutor),
		instanceID: config.InstanceID,
		config:     config,
		logger:     logger,
	}
	return c
}

// RegisterExecutor 注册任务执行器
func (d *DistributedCron) RegisterExecutor(executor TaskExecutor) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.executors[executor.Name()] = executor
}

// AddTask 添加任务
func (d *DistributedCron) AddTask(task *Task) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	if task.ID == "" {
		return fmt.Errorf("task ID cannot be empty")
	}
	if _, exists := d.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}
	d.tasks[task.ID] = task
	if task.Enabled {
		return d.scheduleTask(task)
	}
	return nil
}

// scheduleTask 调度任务
func (d *DistributedCron) scheduleTask(task *Task) error {
	entryID, err := d.cron.AddFunc(task.Expression, d.createTaskFunc(task))
	if err != nil {
		return fmt.Errorf("failed to schedule task %s: %w", task.ID, err)
	}
	task.CronJobID = int(entryID)
	d.logger.Infof("Task %s scheduled with entry ID: %d", task.ID, entryID)
	return nil
}

// createTaskFunc 创建任务执行函数
func (d *DistributedCron) createTaskFunc(task *Task) func() {
	return func() {
		ctx := context.Background()

		// 检查任务是否已经在运行
		if _, running := d.runningTasks.Load(task.ID); running {
			d.logger.Infof("Task %s is already running, skip", task.ID)
			return
		}

		// 标记任务为运行中
		d.runningTasks.Store(task.ID, true)
		defer d.runningTasks.Delete(task.ID)

		// 分布式环境下获取锁
		if d.config.EnableDistributed {
			lockKey := fmt.Sprintf("task_lock:%s", task.ID)
			acquired, err := d.locker.Acquire(ctx, lockKey, d.config.LockTTL)
			if err != nil {
				d.logger.Errorf("Failed to acquire lock for task %s: %v", task.ID, err)
				return
			}

			if !acquired {
				d.logger.Infof("Task %s lock acquired by other instance, skip", task.ID)
				return
			}

			// 确保锁被释放
			defer func() {
				if err := d.locker.Release(ctx, lockKey); err != nil {
					d.logger.Errorf("Failed to release lock for task %s: %v", task.ID, err)
				}
			}()

			// 启动锁刷新
			stopRefresh := make(chan struct{})
			go d.refreshLock(ctx, lockKey, stopRefresh)
			defer close(stopRefresh)
		}

		// 执行任务
		if err := d.executeTask(ctx, task); err != nil {
			d.logger.Errorf("Failed to execute task %s: %v", task.ID, err)
		}
	}
}

// refreshLock 刷新锁的TTL
func (d *DistributedCron) refreshLock(ctx context.Context, lockKey string, stop chan struct{}) {
	ticker := time.NewTicker(d.config.LockRefreshInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if err := d.locker.Refresh(ctx, lockKey, d.config.LockTTL); err != nil {
				d.logger.Errorf("Failed to refresh lock %s: %v", lockKey, err)
				return
			}
		case <-stop:
			return
		}
	}
}

// executeTask 执行任务
func (d *DistributedCron) executeTask(ctx context.Context, task *Task) error {
	executor, exists := d.executors[task.Executor]
	if !exists {
		return fmt.Errorf("executor %s not found for task %s", task.Executor, task.ID)
	}

	startTime := time.Now()
	d.logger.Infof("Starting task %s at %s", task.ID, startTime.Format(time.RFC3339))

	if err := executor.Execute(ctx, task); err != nil {
		return fmt.Errorf("task %s execution failed: %w", task.ID, err)
	}

	duration := time.Since(startTime)
	d.logger.Infof("Completed task %s in %v", task.ID, duration)
	return nil
}

// Start 启动定时任务
func (d *DistributedCron) Start() {
	d.cron.Start()
	d.logger.Infof("Distributed cron started with instance ID: %s", d.instanceID)
}

// Stop 停止定时任务
func (d *DistributedCron) Stop() {
	d.cron.Stop()
	d.logger.Infof("Distributed cron stopped")
}

// GetTasks 获取所有任务
func (d *DistributedCron) GetTasks() []*Task {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	tasks := make([]*Task, 0, len(d.tasks))
	for _, task := range d.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// EnableTask 启用任务
func (d *DistributedCron) EnableTask(taskID string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	task, exists := d.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}
	if task.Enabled {
		return nil
	}
	task.Enabled = true
	return d.scheduleTask(task)
}

// DisableTask 禁用任务
func (d *DistributedCron) DisableTask(taskID string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	task, exists := d.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}
	task.Enabled = false
	d.cron.Remove(cron.EntryID(task.CronJobID))
	return nil
}

// DeleteTask 删除任务
func (d *DistributedCron) DeleteTask(taskID string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	task, exists := d.tasks[taskID]
	if !exists {
		return fmt.Errorf("task %s not found", taskID)
	}
	delete(d.tasks, taskID)
	d.cron.Remove(cron.EntryID(task.CronJobID))
	return nil
}
