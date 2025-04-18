package handlers

import (
	"context"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var _ ActionInfo = (*MessageActionInfo)(nil)

type ActionInfo interface {
}

type MessageActionInfo struct {
	HandlerType HandlerType
	MsgType     string
	MsgId       *string
	ChatId      *string
	QParsed     string
	FileKey     string
	ImageKey    string
	ImageKeys   []string // post 消息卡片中的图片组
	SessionId   *string
	Mention     []*larkim.MentionEvent
}

type ActionData struct {
	ctx     context.Context
	handler *MessageHandler
	info    ActionInfo
}

type Action interface {
	Execute(a *ActionData) error
}

// 责任链
func chain(data *ActionData, actions ...Action) error {
	for _, act := range actions {
		if err := act.Execute(data); err != nil {
			return err
		}
	}
	return nil
}
