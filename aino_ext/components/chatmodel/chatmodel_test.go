package chatmodel_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	_ = godotenv.Load("../../../.env")
	fmt.Println(os.Getenv("ENV"))
}

func TestArkChatModel(t *testing.T) {
	llmBaseUrl := os.Getenv("LLM_BASE_URL")
	llmApiKey := os.Getenv("LLM_API_KEY")
	llmModel := os.Getenv("LLM_MODEL")
}
