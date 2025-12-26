package agent

import (
	"bufio"
	"context"
	"embed"
	"errors"
	"io"
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/hertz-contrib/sse"
	"github.com/zhangga/aino/pkg/mempkg"
	"github.com/zhangga/aino/pkg/utils"
	logger "github.com/zhangga/aino/pkg/zlog"
)

//go:embed web
var webContent embed.FS

type ChatRequest struct {
	Id      string `json:"id"`
	Message string `json:"message"`
}

func BindRoutes(r *route.RouterGroup) error {
	if err := Init(); err != nil {
		return err
	}

	// API 路由
	r.GET("/api/chat", HandleChat)
	r.GET("/api/log", HandleLog)
	r.GET("/api/history", HandleHistory)
	r.DELETE("/api/history", HandleDeleteHistory)

	// 静态文件服务
	r.GET("/", func(ctx context.Context, c *app.RequestContext) {
		content, err := webContent.ReadFile("web/index.html")
		if err != nil {
			c.String(consts.StatusNotFound, "File not found")
			return
		}
		c.Header("Content-Type", "text/html")
		c.Write(content)
	})

	r.GET("/:file", func(ctx context.Context, c *app.RequestContext) {
		file := c.Param("file")
		content, err := webContent.ReadFile("web/" + file)
		if err != nil {
			c.String(consts.StatusNotFound, "File not found")
			return
		}

		contentType := mime.TypeByExtension(filepath.Ext(file))
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		c.Header("Content-Type", contentType)
		c.Write(content)
	})

	return nil
}

// HandleChat 聊天处理
func HandleChat(ctx context.Context, c *app.RequestContext) {
	id := c.Query("id")
	message := c.Query("message")

	if len(id) == 0 || len(message) == 0 {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"status": "error",
			"error":  "id and message are required",
		})
		return
	}

	logger.Infof("[Chat] Starting chat with id: %s, Message: %s", id, message)

	sr, err := RunAgent(ctx, id, message)
	if err != nil {
		logger.Errorf("[Chat] Error id: %s running agent: %v", id, err)
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	s := sse.NewStream(c)
	defer func() {
		sr.Close()
		_ = c.Flush()
		logger.Infof("[Chat] Finished chat with id: %s", id)
	}()

outer:
	for {
		select {
		case <-ctx.Done():
			logger.Infof("[Chat] Context done id: %s: %v", id, ctx.Err())
			return
		default:
			msg, err := sr.Recv()
			if errors.Is(err, io.EOF) {
				logger.Infof("[Chat] EOF received id: %s", id)
				break outer
			}
			if err != nil {
				logger.Errorf("[Chat] Error id: %s receiving message: %v", id, err)
				break outer
			}

			err = s.Publish(&sse.Event{
				Data: utils.StringToBytes(msg.Content),
			})
			if err != nil {
				logger.Errorf("[Chat] Error id: %s publishing message: %v", id, err)
				break outer
			}
		}
	}
}

func HandleHistory(ctx context.Context, c *app.RequestContext) {
	// query: id => get history, none => list all
	id := c.Query("id")

	if id == "" {
		ids := mempkg.GetDefaultMemory().ListConversations()

		c.JSON(consts.StatusOK, map[string]interface{}{
			"ids": ids,
		})
		return
	}

	conversation := mempkg.GetDefaultMemory().GetConversation(id, false)
	if conversation == nil {
		c.JSON(consts.StatusNotFound, map[string]string{
			"error": "conversation not found",
		})
		return
	}

	c.JSON(consts.StatusOK, map[string]interface{}{
		"conversation": conversation,
	})
}

func HandleDeleteHistory(ctx context.Context, c *app.RequestContext) {
	id := c.Query("id")
	if id == "" {
		c.JSON(consts.StatusBadRequest, map[string]string{
			"error": "missing id parameter",
		})
		return
	}

	if err := mempkg.GetDefaultMemory().DeleteConversation(id); err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	c.JSON(consts.StatusOK, map[string]string{
		"status": "success",
	})
}

// HandleLog 实时日志处理
func HandleLog(ctx context.Context, c *app.RequestContext) {
	file, err := os.Open("logs/einoagent_detail.log")
	if err != nil {
		c.JSON(consts.StatusInternalServerError, map[string]string{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}
	defer file.Close()

	// Create a new SSE stream
	s := sse.NewStream(c)
	defer c.Flush()

	// Seek to the end of the file
	_, err = file.Seek(0, io.SeekEnd)
	if err != nil {
		logger.Infof("error seeking file: %v", err)
		return
	}

	// Use a goroutine to continuously read new log lines and send them to the client
	go func() {
		reader := bufio.NewReader(file)
		for {
			line, err := reader.ReadString('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				logger.Infof("error reading log file: %v", err)
				break
			}

			// If we got a line, publish it
			if len(line) > 0 {
				err = s.Publish(&sse.Event{
					Data: utils.StringToBytes(line),
				})
				if err != nil {
					logger.Infof("error publishing log line: %v", err)
					break
				}
			}

			// If we reached EOF, wait a bit before trying again
			if errors.Is(err, io.EOF) {
				select {
				case <-ctx.Done():
					logger.Infof("context done: %v", ctx.Err())
					return
				case <-time.After(5 * time.Second):
				}
			}
		}
	}()

	// Keep the connection open
	<-ctx.Done()
}
