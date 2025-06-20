package service

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/zhangga/aino/internal/conf"
	"github.com/zhangga/aino/internal/service/handlers"
	"github.com/zhangga/aino/internal/service/sctx"
	"github.com/zhangga/aino/pkg/logger"
	"runtime/debug"
)

var service *Service

type Service struct {
	ctx            context.Context
	srvConf        *conf.ServiceConfig
	taskChan       chan Task
	messageHandler handlers.MessageHandlerInterface
	larkClient     *lark.Client
}

func RunService(ctx context.Context, srvConf *conf.ServiceConfig, larkConf *conf.LarkConfig) error {
	service = &Service{
		ctx:            ctx,
		srvConf:        srvConf,
		taskChan:       make(chan Task, 1024),
		messageHandler: handlers.NewMessageHandler(),
		larkClient:     lark.NewClient(larkConf.AppID, larkConf.AppSecret),
	}
	return service.Start()
}

func (srv *Service) Start() error {
	errAccepted := make(chan error)
	// 启动http服务
	go func() {
		r := gin.Default()
		ginHandlers(r)
		logger.Infof("http server started: http://localhost:%d/ping\n", srv.srvConf.HttpPort)
		if err := r.Run(fmt.Sprintf(":%d", srv.srvConf.HttpPort)); err != nil {
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

	ctx := sctx.WithContext(srv.ctx).WithServiceConfig(srv.srvConf).WithLarkClient(srv.larkClient)

	actionInfo, err := task.AsActionInfo()
	if err != nil {
		logger.Errorf("[Task] taskId=%d, taskType=%d, AsData error: %v", task.Id(), task.Type(), err)
		return
	}

	switch task.Type() {
	case TaskTypeLark:
		if err = srv.messageHandler.MsgReceivedHandler(ctx, actionInfo); err != nil {
			logger.Errorf("[Task] taskId=%d, taskType=%d, MsgReceivedHandler error: %v", task.Id(), task.Type(), err)
			return
		}
	case TaskTypeCard:
		if err = srv.messageHandler.CardHandler(ctx, actionInfo); err != nil {
			logger.Errorf("[Task] taskId=%d, taskType=%d, CardActionHandler error: %v", task.Id(), task.Type(), err)
			return
		}
	default:
		logger.Errorf("[Task] taskId=%d, taskType=%d, not support", task.Id(), task.Type())
		return
	}
}
