package dcron

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestDcron(t *testing.T) {
	// 根据环境选择锁实现
	var locker DistributedLocker
	var enableDistributed bool
	_ = os.Setenv("REDIS_URL", "127.0.0.1:6379")
	_ = os.Setenv("REDIS_PASSWORD", "")
	if os.Getenv("REDIS_URL") != "" {
		redisClient := redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       15,
		})
		locker = NewRedisLocker("dcron", redisClient)
		enableDistributed = true
		log.Println("Running in distributed mode with Redis locker")
	}

	// 创建配置
	config := &Config{
		InstanceID:          getInstanceID(),
		LockTTL:             30 * time.Second,
		LockRefreshInterval: 10 * time.Second,
		EnableDistributed:   enableDistributed,
	}

	// 创建分布式定时任务管理器
	cronManager := NewDistributedCron(config, locker, nil)

	// 注册执行器
	cronManager.RegisterExecutor(NewSimplePrintExecutor("print"))
	cronManager.RegisterExecutor(NewHTTPRequestExecutor("http"))

	// 添加任务
	tasks := []*Task{
		{
			ID:          "task-1",
			Name:        "每10秒打印任务",
			Expression:  "*/10 * * * * *", // 每10秒
			Executor:    "print",
			Parameters:  map[string]any{"message": "Hello from distributed cron!"},
			Description: "简单的打印任务",
			Enabled:     true,
		},
		{
			ID:          "task-2",
			Name:        "每30秒HTTP请求",
			Expression:  "*/30 * * * * *", // 每30秒
			Executor:    "http",
			Parameters:  map[string]any{"url": "https://httpbin.org/get"},
			Description: "HTTP请求任务",
			Enabled:     true,
		},
	}

	for _, task := range tasks {
		if err := cronManager.AddTask(task); err != nil {
			log.Printf("Failed to add task %s: %v", task.ID, err)
		}
	}

	// 启动定时任务
	cronManager.Start()

	// 等待中断信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// 优雅关闭
	log.Println("Shutting down...")
	cronManager.Stop()
	log.Println("Shutdown completed")
}
