package chatmodel_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load("../../../.env")
	fmt.Println(os.Getenv("ENV"))
}

func TestArkChatModel(t *testing.T) {
	arkApiKey := os.Getenv("ARK_API_KEY")
	arkModel := os.Getenv("ARK_CHAT_MODEL")

	config := &ark.ChatModelConfig{
		Model:  arkModel,
		APIKey: arkApiKey,
	}

	ctx := context.Background()
	cm, err := ark.NewChatModel(ctx, config)
	if err != nil {
		t.Fatalf("failed to create Ark chat model: %v", err)
	}

	input := []*schema.Message{
		&schema.Message{
			Role:    schema.System,
			Content: "You are a helpful assistant that translates Chinese to English.",
		},
		&schema.Message{
			Role:    schema.User,
			Content: "做一个详细的自我介绍。",
		},
	}
	outStream, err := cm.Stream(ctx, input)
	if err != nil {
		t.Fatalf("failed to call Ark chat model: %v", err)
	}

	ticker := time.Tick(500 * time.Millisecond)
	for {
		select {
		case <-ticker:
			msg, err := outStream.Recv()
			if errors.Is(err, io.EOF) {
				t.Log("stream ended")
				return
			}
			if err != nil {
				t.Fatalf("failed to receive message from stream: %v", err)
			}
			if msg == nil {
				t.Log("stream ended")
				return
			}
			if msg.Content != "" {
				t.Logf("received message: %s", msg.Content)
			}
		}
	}
}
