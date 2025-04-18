package handlers

import (
	"context"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

type MessageHandlerInterface interface {
	MsgReceivedHandler(ctx context.Context, data ActionInfo) error
	CardHandler(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error)
}

type HandlerType string

const (
	GroupHandler   = "group"
	UserHandler    = "personal"
	UnknownHandler = "unknown"
)
