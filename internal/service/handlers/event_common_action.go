package handlers

import "errors"

var (
	_ Action = (*ProcessedUniqueAction)(nil)
)

var (
	ErrProcessedUniqueAction = errors.New("processed unique action")
)

type ProcessedUniqueAction struct { //消息唯一性
}

func (ProcessedUniqueAction) Execute(data *ActionData) error {
	if data.handler.msgCache.IfProcessed(data.info.Id()) {
		return ErrProcessedUniqueAction
	}
	data.handler.msgCache.TagProcessed(data.info.Id())
	return nil
}
