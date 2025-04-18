package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zhangga/aino/internal/service/cache"
	"github.com/zhangga/aino/internal/service/handlers"
	"github.com/zhangga/aino/pkg/logger"
	"runtime/debug"
)

var service *Service

type Service struct {
	ctx            context.Context
	httpPort       int
	taskChan       chan Task
	messageHandler handlers.MessageHandlerInterface
}

func RunService(ctx context.Context, httpPort int) error {
	msgCache := cache.NewMsgCache()
	service = &Service{
		ctx:            ctx,
		httpPort:       httpPort,
		taskChan:       make(chan Task, 1024),
		messageHandler: handlers.NewMessageHandler(msgCache),
	}
	return service.Start()
}

func (srv *Service) Start() error {
	errAccepted := make(chan error)
	// 启动http服务
	go func() {
		r := gin.Default()
		ginHandlers(r)
		logger.Infof("http server started: http://localhost:%d/ping\n", srv.httpPort)
		if err := r.Run(fmt.Sprintf(":%d", srv.httpPort)); err != nil {
			errAccepted <- err
		}
	}()

	go func() {
		for {
			select {
			case task := <-srv.taskChan:
				logger.Debugf("[Service] handle taskId=%d", task.Id())
				// 启动一个goroutine执行任务
				go srv.SafeHandle(task)
			case <-srv.ctx.Done(): // 主程序退出
				return
			case err := <-errAccepted:
				if err != nil {
					logger.Errorf("[Service] http server start failed: %v", err)
					return
				}
			}
		}
	}()
	return nil
}

func (srv *Service) AddTask(task Task) error {
	logger.Debugf("[Service] add taskId=%d, taskType=%d", task.Id(), task.Type())
	srv.taskChan <- task
	return nil
}

func (srv *Service) SafeHandle(task Task) {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("[Task] taskId=%d, taskType=%d Run() panic: %v\n%s", task.Id(), task.Type(), err, debug.Stack())
		}
	}()

	actionInfo, err := task.AsActionInfo()
	if err != nil {
		logger.Errorf("[Task] taskId=%d, taskType=%d, AsData error: %v", task.Id(), task.Type(), err)
		return
	}

	switch task.Type() {
	case TaskTypeLark:
		if err = srv.messageHandler.MsgReceivedHandler(srv.ctx, actionInfo); err != nil {
			logger.Errorf("[Task] taskId=%d, taskType=%d, MsgReceivedHandler error: %v", task.Id(), task.Type(), err)
			return
		}
	default:
		logger.Errorf("[Task] taskId=%d, taskType=%d, not support", task.Id(), task.Type())
		return
	}
}
