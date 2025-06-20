package handlers

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/schema"
	"github.com/zhangga/aino/internal/eino"
	"github.com/zhangga/aino/internal/service/sctx"
	"github.com/zhangga/aino/internal/tools"
	"github.com/zhangga/aino/pkg/logger"
	"io"
	"time"
)

var (
	ErrorResponseStream = errors.New("response stream")
	ErrorUpdateCard     = errors.New("update card")
)

var _ Action = (*StreamMessageAction)(nil)

type StreamMessageAction struct { /*消息*/
}

func (s *StreamMessageAction) Execute(a *ActionData) error {
	srvConf := sctx.GetServiceConfig(a.ctx)
	if !srvConf.StreamMode {
		return nil
	}

	// 获取会话历史消息
	msg := a.handler.sessionCache.GetMsg(a.info.GetSessionId())
	// 系统提示词
	var systemPrompt string
	if systemMsg, index := findSystemRole(msg); index != -1 {
		systemPrompt = systemMsg.Content
		// 删除原来的系统提示词
		msg = append(msg[:index], msg[index+1:]...)
		msg = msg[:len(msg)-1]
	} else {
		systemPrompt = getDefaultSystemPrompt()
	}
	// 添加新消息
	msg = append(msg, &schema.Message{
		Role:    schema.User,
		Content: a.info.GetContent(),
	})

	// if new topic
	var ifNewTopic bool
	if len(msg) <= 2 {
		ifNewTopic = true
	} else {
		ifNewTopic = false
	}

	// 发送消息卡片，正在处理中
	cardId, err2 := sendOnProcess(a, ifNewTopic)
	if err2 != nil {
		return err2
	}

	answer := ""
	chatResponseStream := make(chan string)
	done := make(chan struct{}) // 添加 done 信号，保证 goroutine 正确退出
	noContentTimeout := time.AfterFunc(1000*time.Second, func() {
		logger.Error("no content timeout")
		close(done)
		err := updateFinalCard(a.ctx, "请求超时", cardId, ifNewTopic)
		if err != nil {
			return
		}
		return
	})
	defer noContentTimeout.Stop()

	// 请求大模型
	go func() {
		defer func() {
			if err := recover(); err != nil {
				_ = updateFinalCard(a.ctx, fmt.Sprintf("聊天失败: %s", err), cardId, ifNewTopic)
			}
		}()

		ragent, err := eino.NewAgentByType(a.ctx, systemPrompt, tools.AllCreator...)
		if err != nil {
			panic("创建Eino Agent失败")
		}
		output, err := ragent.Stream(a.ctx, msg)
		if err != nil {
			panic("发送Stream消息失败")
		}

		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop() // 注意在函数结束时停止 ticker

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				frame, err := output.Recv()
				if errors.Is(err, io.EOF) {
					close(done)
					return
				}
				if err != nil {
					panic(err)
				}
				chatResponseStream <- frame.Content
			}
		}

		//aiMode := a.handler.sessionCache.GetAIMode(a.info.GetSessionId())
		////fmt.Println("msg: ", msg)
		////fmt.Println("aiMode: ", aiMode)
		//if err := a.handler.gpt.StreamChat(a.ctx, msg, aiMode,
		//	chatResponseStream); err != nil {
		//	err := updateFinalCard(a.ctx, "聊天失败", cardId, ifNewTopic)
		//	if err != nil {
		//		return
		//	}
		//	close(done) // 关闭 done 信号
		//}
		//
		//close(done) // 关闭 done 信号
	}()

	for {
		select {
		case res, ok := <-chatResponseStream:
			if !ok {
				return ErrorResponseStream
			}
			noContentTimeout.Stop()
			answer += res
			if err := updateTextCard(a.ctx, answer, cardId, ifNewTopic); err != nil {
				return err
			}
		case <-done: // 添加 done 信号的处理
			err := updateFinalCard(a.ctx, answer, cardId, ifNewTopic)
			if err != nil {
				return ErrorUpdateCard
			}

			msg = append(msg, &schema.Message{
				Role: "assistant", Content: answer,
			})
			a.handler.sessionCache.SetMsg(a.info.GetSessionId(), msg)
			close(chatResponseStream)

			jsonByteArray, err := sonic.Marshal(msg)
			if err != nil {
				logger.Errorf("sonic marshal err %v", err)
				return nil
			}
			logger.Debugf("StreamMessageAction.Execute() result: %s", string(jsonByteArray))
			return nil
		}
	}
}

func setDefaultPrompt(msg []*schema.Message) []*schema.Message {
	if _, index := findSystemRole(msg); index == -1 {
		msg = append(msg, &schema.Message{
			Role:    schema.System,
			Content: getDefaultSystemPrompt(),
		})
	}
	return msg
}

// 判断msg中的是否包含system role
func findSystemRole(msg []*schema.Message) (*schema.Message, int) {
	for i, m := range msg {
		if m.Role == schema.System {
			return m, i
		}
	}
	return nil, -1
}

func getDefaultSystemPrompt() string {
	//return "You are ChatGPT, a large language model trained by OpenAI. Answer in user's language as concisely as possible. Knowledge cutoff: 20230601 Current date" + time.Now().Format("20060102")
	return "你是一个人工智能助手，可以解决日常碰到的任何问题。\n1. 我的输入是飞书消息\n2. 你根据我的飞书消息内容，并充分利用现有工具来解决问题\n3. 在你认为任务真正完成，或已无法完成的情况下才停止"
}

func sendOnProcess(a *ActionData, ifNewTopic bool) (string, error) {
	// send 正在处理中
	cardId, err := sendOnProcessCard(a.ctx, a.info.GetSessionId(), a.info.GetMsgId(), ifNewTopic)
	if err != nil {
		return "", err
	}
	return cardId, nil

}
