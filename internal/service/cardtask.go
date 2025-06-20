package service

import (
	"context"
	"github.com/bytedance/sonic"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher/callback"
	"github.com/zhangga/aino/internal/service/handlers"
)

var _ Task = (*CardTask)(nil)

// CardTask 处理卡片消息的任务
type CardTask struct {
	id    uint64
	event *callback.CardActionTriggerEvent
}

func NewCardTask(event *callback.CardActionTriggerEvent) *CardTask {
	task := &CardTask{
		id:    NextID(),
		event: event,
	}
	return task
}

func (t *CardTask) Id() uint64 {
	return t.id
}

func (t *CardTask) Type() TaskType {
	return TaskTypeCard
}

func (t *CardTask) AsActionInfo() (handlers.ActionInfo, error) {
	action := t.event.Event.Action
	var value handlers.CardActionValue
	bs, err := sonic.Marshal(action.Value)
	if err != nil {
		return nil, err
	}
	if err = sonic.Unmarshal(bs, &value); err != nil {
		return nil, err
	}

	info := &handlers.CardActionInfo{
		ChatType:   handlers.ChatUser,
		Value:      value,
		Tag:        action.Tag,
		Option:     action.Option,
		Name:       action.Name,
		InputValue: action.InputValue,
		Checked:    action.Checked,
	}
	return info, nil
}

func (t *CardTask) Run(ctx context.Context) {
	//TODO implement me
	panic("implement me")
}
