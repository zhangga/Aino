package common

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

var _ model.ToolCallingChatModel = (*ChatModelImpl)(nil)

type ChatModelImpl struct {
	client     *http.Client
	baseURL    string
	apiKey     string
	model      string
	retryCount int
	timeout    time.Duration
}

type ChatModelConfig struct {
	BaseURL string
	APIKey  string
	Model   string
}

type ChatModelOptions struct {
	BaseOptions *model.Options
	RetryCount  int
	Timeout     time.Duration
}

func NewChatModel(config *ChatModelConfig) (*ChatModelImpl, error) {
	if config.APIKey == "" {
		return nil, errors.New("APIKey is required")
	}
	return &ChatModelImpl{
		client:  &http.Client{},
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		model:   config.Model,
	}, nil
}

func (m *ChatModelImpl) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// 1. 处理选项
	baseOpts := model.GetCommonOptions(nil, opts...)
	baseOpts.Model = &m.model
	options := &ChatModelOptions{
		BaseOptions: baseOpts,
		RetryCount:  m.retryCount,
		Timeout:     m.timeout,
	}
	options = model.GetImplSpecificOptions(options, opts...)

	// 2. 开始生成前的回调
	ctx = callbacks.OnStart(ctx, &model.CallbackInput{
		Messages: input,
		Config: &model.Config{
			Model: *options.BaseOptions.Model,
		},
	})

	// 3. 模型推理逻辑
	response, err := m.doGenerate(ctx, input, options)

	// 4. 处理错误
	if err != nil {
		callbacks.OnError(ctx, err)
		return nil, err
	}

	// 5. 结束回调
	callbacks.OnEnd(ctx, &model.CallbackOutput{
		Message: response,
	})
	return response, nil
}

func (m *ChatModelImpl) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	// 1. 处理选项
	baseOpts := model.GetCommonOptions(nil, opts...)
	baseOpts.Model = &m.model
	options := &ChatModelOptions{
		BaseOptions: baseOpts,
		RetryCount:  m.retryCount,
		Timeout:     m.timeout,
	}
	options = model.GetImplSpecificOptions(options, opts...)

	// 2. 开始流式生成前的回调
	ctx = callbacks.OnStart(ctx, &model.CallbackInput{
		Messages: input,
		Config: &model.Config{
			Model: *options.BaseOptions.Model,
		},
	})

	// 3. 创建流式响应
	// Pipe产生一个StreamReader和一个StreamWrite，向StreamWrite中写入可以从StreamReader中读到，二者并发安全。
	// 实现中异步向StreamWrite中写入生成内容，返回StreamReader作为返回值
	// ***StreamReader是一个数据流，仅可读一次，组件自行实现Callback时，既需要通过OnEndWithCallbackOutput向callback传递数据流，也需要向返回一个数据流，需要对数据流进行一次拷贝
	// 考虑到此种情形总是需要拷贝数据流，OnEndWithCallbackOutput函数会在内部拷贝并返回一个未被读取的流
	// 以下代码演示了一种流处理方式，处理方式不唯一
	reader, writer := schema.Pipe[*model.CallbackOutput](1)

	// 4. 启动协程处理流式生成
	go func() {
		defer writer.Close()
		// 模型流式推理逻辑
		m.doStream(ctx, input, options, writer)
	}()

	// 5. 结束回调
	_, nsr := callbacks.OnEndWithStreamOutput(ctx, reader)
	return schema.StreamReaderWithConvert(nsr, func(t *model.CallbackOutput) (*schema.Message, error) {
		return t.Message, nil
	}), nil
}

func (m *ChatModelImpl) WithTools(tools []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	// 实现工具绑定逻辑
	return m, nil
}

func (m *ChatModelImpl) doGenerate(ctx context.Context, input []*schema.Message, options *ChatModelOptions) (*schema.Message, error) {
	return nil, nil
}

func (m *ChatModelImpl) doStream(ctx context.Context, input []*schema.Message, options *ChatModelOptions, sw *schema.StreamWriter[*model.CallbackOutput]) {

}
