package eino_test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/schema"
	"github.com/stretchr/testify/assert"
	"github.com/zhangga/aino/internal/conf"
	"github.com/zhangga/aino/internal/eino"
	"os"
	"testing"
)

func initAgent() {
	// 获取环境变量
	llmConfig := &conf.LLMConfig{
		Model:   os.Getenv("LLM_MODEL"),
		ApiKey:  os.Getenv("LLM_API_KEY"),
		BaseURL: os.Getenv("LLM_BASE_URL"),
	}
	err := eino.InitAgent(context.Background(), llmConfig)
	if err != nil {
		panic(err)
	}
}

func TestAgent(t *testing.T) {
	initAgent()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	persona := `# Character:
你是一个帮助用户推荐餐厅和菜品的助手，根据用户的需要，查询餐厅信息并推荐，查询餐厅的菜品并推荐。
`
	ragent, err := eino.NewAgent(ctx, persona, &ToolQueryRestaurants{})
	assert.Nil(t, err)

	msg, err := ragent.Generate(ctx, []*schema.Message{
		{
			Role:    schema.User,
			Content: "我在北京，给我推荐一些菜，需要有口味辣一点的菜，至少推荐有 2 家餐厅",
		},
	}, agent.WithComposeOptions(compose.WithCallbacks(&eino.LoggerCallback{})))
	assert.Nil(t, err)

	t.Logf("msg: %s", msg)
}

type ToolQueryRestaurants struct {
}

func (t *ToolQueryRestaurants) Info(ctx context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "query_restaurants",
		Desc: "Query restaurants",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"location": {
				Type:     "string",
				Desc:     "The location of the restaurant",
				Required: true,
			},
			"topn": {
				Type: "number",
				Desc: "top n restaurant in some location sorted by score",
			},
		}),
	}, nil
}

type QueryRestaurantsParam struct {
	Location string `json:"location"`
	Topn     int    `json:"topn"`
}

// InvokableRun
// tool 接收的参数和返回都是 string, 就如大模型的 tool call 的返回一样, 因此需要自行处理参数和结果的序列化.
// 返回的 content 会作为 schema.Message 的 content, 一般来说是作为大模型的输入, 因此处理成大模型能更好理解的结构最好.
// 因此，如果是 json 格式，就需要注意 key 和 value 的表意, 不要用 int Enum 代表一个业务含义，比如 `不要用 1 代表 male, 2 代表 female` 这类.
func (t *ToolQueryRestaurants) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// 解析参数
	p := &QueryRestaurantsParam{}
	err := json.Unmarshal([]byte(argumentsInJSON), p)
	if err != nil {
		return "", err
	}
	if p.Topn == 0 {
		p.Topn = 3
	}

	// 请求后端服务
	fmt.Println("=========[ToolQueryRestaurants.ToolQueryRestaurants()] =========")

	return "牛逼餐厅1", nil
}
