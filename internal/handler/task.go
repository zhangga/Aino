package handler

import (
	"context"
	"github.com/zhangga/aino/pkg/logger"
	"runtime/debug"
)

type TaskType int

const (
	TaskTypeLark TaskType = iota
)

type Task interface {
	Id() uint64
	Type() TaskType
	Run(ctx context.Context)
}

func SafeRun(ctx context.Context, task Task) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("[Task] taskId=%d, taskType=%d Run() panic: %v\n%s", task.Id(), task.Type(), err, debug.Stack())
		}
	}()

	task.Run(ctx)
}

func AddTask(task Task) error {
	return service.AddTask(task)
}
