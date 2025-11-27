package dcron

import (
	"context"
	"fmt"
	"log"
	"time"
)

// SimplePrintExecutor 简单的打印执行器
type SimplePrintExecutor struct {
	name string
}

func NewSimplePrintExecutor(name string) *SimplePrintExecutor {
	return &SimplePrintExecutor{name: name}
}

func (s *SimplePrintExecutor) Execute(_ context.Context, task *Task) error {
	message := "Hello from task"
	if msg, ok := task.Parameters["message"]; ok {
		message = msg.(string)
	}
	log.Printf("[%s] %s: %s", s.name, task.Name, message)
	time.Sleep(2 * time.Second) // 模拟任务执行时间
	return nil
}

func (s *SimplePrintExecutor) Name() string {
	return s.name
}

// HTTPRequestExecutor HTTP请求执行器
type HTTPRequestExecutor struct {
	name string
}

func NewHTTPRequestExecutor(name string) *HTTPRequestExecutor {
	return &HTTPRequestExecutor{name: name}
}

func (h *HTTPRequestExecutor) Execute(_ context.Context, task *Task) error {
	url, ok := task.Parameters["url"].(string)
	if !ok {
		return fmt.Errorf("url parameter is required")
	}

	log.Printf("Making HTTP request to: %s", url)
	// 这里实现实际的HTTP请求逻辑
	time.Sleep(1 * time.Second)
	log.Printf("HTTP request to %s completed", url)
	return nil
}

func (h *HTTPRequestExecutor) Name() string {
	return h.name
}
