package handlers

import (
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

var (
	_ ActionInfo = (*MessageActionInfo)(nil)
	_ ActionInfo = (*CardActionInfo)(nil)
)

type ActionInfo interface {
	GetChatType() ChatType
	GetMsgId() string
	GetSessionId() string
	GetMsgType() string
	GetContent() string
}

type MessageActionInfo struct {
	ChatType  ChatType
	MsgType   string
	MsgId     string
	ChatId    string
	Content   string
	FileKey   string
	ImageKey  string
	ImageKeys []string // post 消息卡片中的图片组
	SessionId string
	Mention   []*larkim.MentionEvent
}

func (m *MessageActionInfo) GetChatType() ChatType {
	return m.ChatType
}

func (m *MessageActionInfo) GetMsgId() string {
	return m.MsgId
}

func (m *MessageActionInfo) GetSessionId() string {
	return m.SessionId
}

func (m *MessageActionInfo) GetMsgType() string {
	return m.MsgType
}

func (m *MessageActionInfo) GetContent() string {
	return m.Content
}

type CardActionValue struct {
	Kind      CardKind `json:"kind"`
	MsgId     string   `json:"msgId"`
	SessionId string   `json:"sessionId"`
	Value     string   `json:"value"`
}

type CardActionInfo struct {
	ChatType   ChatType
	Value      CardActionValue
	Tag        string
	Option     string
	Name       string
	InputValue string
	Checked    bool
}

func (c *CardActionInfo) GetChatType() ChatType {
	return c.ChatType
}

func (c *CardActionInfo) GetMsgId() string {
	return c.Value.MsgId
}

func (c *CardActionInfo) GetSessionId() string {
	return c.Value.SessionId
}

func (c *CardActionInfo) GetMsgType() string {
	return string(c.Value.Kind)
}

func (c *CardActionInfo) GetContent() string {
	return c.Value.Value
}
