package handlers

import (
	"context"
)

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
