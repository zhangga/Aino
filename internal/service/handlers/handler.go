package handlers

import (
	"context"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
	"github.com/zhangga/aino/pkg/logger"
)

var _ MessageHandlerInterface = (*MessageHandler)(nil)

type MessageHandler struct {
	msgCache MsgCacheInterface
}

func NewMessageHandler() MessageHandlerInterface {
	return &MessageHandler{}
}

func (m MessageHandler) MsgReceivedHandler(ctx context.Context, info ActionInfo) error {
	logger.Debugf("处理消息: %v", info)
	data := &ActionData{
		ctx:     ctx,
		handler: &m,
		info:    info,
	}
	actions := []Action{
		&ProcessedUniqueAction{}, //避免重复处理
	}
	return chain(data, actions...)
}

func (m MessageHandler) CardHandler(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}
