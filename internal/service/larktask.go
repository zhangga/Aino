package service

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/zhangga/aino/internal/service/handlers"
	"regexp"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/schema"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/zhangga/aino/internal/eino"
	"github.com/zhangga/aino/internal/tools"
	"github.com/zhangga/aino/pkg/logger"
)

//go:embed larkprompt.txt
var larkPrompt string

var larkTools = []tools.Creator{
	tools.NewToolLarkGetMsg,
	// tools.NewToolLarkSendMsg,
}

var _ Task = (*LarkTask)(nil)

// LarkTask 处理lark消息的任务
type LarkTask struct {
	id  uint64
	Msg *larkim.P2MessageReceiveV1
}

func NewLarkTask(event *larkim.P2MessageReceiveV1) *LarkTask {
	task := &LarkTask{
		id:  NextID(),
		Msg: event,
	}
	return task
}

func (t *LarkTask) Id() uint64 {
	return t.id
}

func (t *LarkTask) Type() TaskType {
	return TaskTypeLark
}

func (t *LarkTask) AsActionInfo() (handlers.ActionInfo, error) {
	handlerType, err := judgeChatType(t.Msg)
	if err != nil {
		return nil, err
	}
	msgType, err := judgeMsgType(t.Msg)
	if err != nil {
		return nil, err
	}

	content := t.Msg.Event.Message.Content
	msgId := t.Msg.Event.Message.MessageId
	rootId := t.Msg.Event.Message.RootId
	chatId := t.Msg.Event.Message.ChatId
	mention := t.Msg.Event.Message.Mentions

	sessionId := rootId
	if sessionId == nil || *sessionId == "" {
		sessionId = msgId
	}

	data := &handlers.MessageActionInfo{
		HandlerType: handlerType,
		MsgType:     msgType,
		MsgId:       msgId,
		ChatId:      chatId,
		SessionId:   sessionId,
		Mention:     mention,
		QParsed:     strings.Trim(parseContent(*content, msgType), " "),
		FileKey:     parseFileKey(*content),
		ImageKey:    parseImageKey(*content),
		ImageKeys:   parsePostImageKeys(*content),
	}
	return data, nil
}

func judgeChatType(msg *larkim.P2MessageReceiveV1) (handlers.HandlerType, error) {
	chatType := msg.Event.Message.ChatType
	switch *chatType {
	case "group":
		return handlers.GroupHandler, nil
	case "p2p":
		return handlers.UserHandler, nil
	default:
		return "", fmt.Errorf("unknow chat type: %s", *chatType)
	}
}

func judgeMsgType(msg *larkim.P2MessageReceiveV1) (string, error) {
	msgType := msg.Event.Message.MessageType
	switch *msgType {
	case "text", "image", "audio", "post":
		return *msgType, nil
	default:
		return "", fmt.Errorf("unknown message type: %v", *msgType)
	}
}

func parseContent(content, msgType string) string {
	//"{\"text\":\"@_user_1  hahaha\"}",
	//only get text content hahaha
	if msgType == "post" {
		return parsePostContent(content)
	}

	var contentMap map[string]interface{}
	err := sonic.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
	}
	if contentMap["text"] == nil {
		return ""
	}
	text := contentMap["text"].(string)
	return msgFilter(text)
}

// Parse rich text json to text
func parsePostContent(content string) string {
	var contentMap map[string]interface{}
	err := sonic.Unmarshal([]byte(content), &contentMap)

	if err != nil {
		fmt.Println(err)
	}

	if contentMap["content"] == nil {
		return ""
	}
	var text string
	// deal with title
	if contentMap["title"] != nil && contentMap["title"] != "" {
		text += contentMap["title"].(string) + "\n"
	}
	// deal with content
	contentList := contentMap["content"].([]interface{})
	for _, v := range contentList {
		for _, v1 := range v.([]interface{}) {
			if v1.(map[string]interface{})["tag"] == "text" {
				text += v1.(map[string]interface{})["text"].(string)
			}
		}
		// add new line
		text += "\n"
	}
	return msgFilter(text)
}

func parseFileKey(content string) string {
	var contentMap map[string]interface{}
	err := sonic.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if contentMap["file_key"] == nil {
		return ""
	}
	fileKey := contentMap["file_key"].(string)
	return fileKey
}

func parseImageKey(content string) string {
	var contentMap map[string]interface{}
	err := sonic.Unmarshal([]byte(content), &contentMap)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if contentMap["image_key"] == nil {
		return ""
	}
	imageKey := contentMap["image_key"].(string)
	return imageKey
}

func parsePostImageKeys(content string) []string {
	var contentMap map[string]interface{}
	err := sonic.Unmarshal([]byte(content), &contentMap)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	var imageKeys []string

	if contentMap["content"] == nil {
		return imageKeys
	}

	contentList := contentMap["content"].([]interface{})
	for _, v := range contentList {
		for _, v1 := range v.([]interface{}) {
			if v1.(map[string]interface{})["tag"] == "img" {
				imageKeys = append(imageKeys, v1.(map[string]interface{})["image_key"].(string))
			}
		}
	}

	return imageKeys
}

// func sendCard
func msgFilter(msg string) string {
	//replace @到下一个非空的字段 为 ''
	regex := regexp.MustCompile(`@[^ ]*`)
	return regex.ReplaceAllString(msg, "")
}

func (t *LarkTask) Run(ctx context.Context) {
	ragent, err := eino.NewAgentByType(ctx, larkPrompt, larkTools...)
	if err != nil {
		panic(err)
	}

	message, err := sonic.Marshal(t.Msg.Event)
	if err != nil {
		panic(err)
	}
	result, err := ragent.Generate(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: string(message),
		},
	}, agent.WithComposeOptions(compose.WithCallbacks(&eino.LoggerCallback{})))
	if err != nil {
		panic(err)
	}

	logger.Infof("Lark taskId=%d, result=%v", t.id, result)

	tools.NewToolLarkSendMsg().(*tools.ToolLarkSendMsg).SendMessage(ctx, *t.Msg.Event.Message.ChatId, string(result.Content))
}
