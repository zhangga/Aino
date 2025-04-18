package service

import (
	"github.com/gin-gonic/gin"
)

func init() {
	registerPostHandler("/webhook/event", handleWebhookEvent)
	registerPostHandler("/webhook/card", handleWebhookCard)
}

func handleWebhookEvent(c *gin.Context) {
	//sdkginext.NewEventHandlerFunc(eventHandler)
}

func handleWebhookCard(c *gin.Context) {
	//sdkginext.NewCardActionHandlerFunc(cardHandler)
}
