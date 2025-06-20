package handlers

import (
	"context"
)

type MessageHandlerInterface interface {
	MsgReceivedHandler(ctx context.Context, data ActionInfo) error
	CardHandler(ctx context.Context, data ActionInfo) error
}

type ChatType string

const (
	ChatGroup   ChatType = "group"
	ChatUser    ChatType = "personal"
	ChatUnknown ChatType = "unknown"
)
