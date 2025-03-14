package handler

import (
	"context"
	"github.com/zhangga/aino/pkg/logger"
)

var service *Service

type Service struct {
	ctx      context.Context
	taskChan chan Task
}

func RunService(ctx context.Context) {
	service = &Service{
		ctx:      ctx,
		taskChan: make(chan Task, 1024),
	}
	service.Start()
}

func (srv *Service) Start() {
	go func() {
		for {
			select {
			case task := <-srv.taskChan:
				logger.Debugf("[Service] handle taskId=%d", task.Id())
				SafeRun(srv.ctx, task)
			case <-srv.ctx.Done(): // 主程序退出
				return
			}
		}
	}()
}

func (srv *Service) AddTask(task Task) error {
	logger.Debugf("[Service] add taskId=%d, taskType=%d", task.Id(), task.Type())
	srv.taskChan <- task
	return nil
}
