package handlers

import larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"

var _ ActionInfo = (*MessageActionInfo)(nil)

type ActionInfo interface {
	Type() HandlerType
	Id() string
}

type MessageActionInfo struct {
	HandlerType HandlerType
	MsgType     string
	MsgId       string
	ChatId      string
	QParsed     string
	FileKey     string
	ImageKey    string
	ImageKeys   []string // post 消息卡片中的图片组
	SessionId   string
	Mention     []*larkim.MentionEvent
}

func (m *MessageActionInfo) Type() HandlerType {
	return m.HandlerType
}

func (m *MessageActionInfo) Id() string {
	return m.MsgId
}
