package handler

import (
	"encoding/json"
	"net/http"
	"github.com/gin-gonic/gin"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/service/event/v1"
)

type LarkHandler struct {
	client *lark.Client
}

func NewLarkHandler(appID, appSecret string) *LarkHandler {
	return &LarkHandler{
		client: lark.NewClient(appID, appSecret),
	}
}

func (h *LarkHandler) SetupRoutes(router *gin.Engine) {
	router.GET("/lark/callback", h.handleVerification)
	router.POST("/lark/callback", h.handleEvent)
}

func (h *LarkHandler) handleVerification(c *gin.Context) {
	var req struct {
		Challenge string `json:"challenge"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"challenge": req.Challenge})
}

func (h *LarkHandler) handleEvent(c *gin.Context) {
	var event larkevent.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 处理事件逻辑
	// TODO: 添加业务处理逻辑

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}