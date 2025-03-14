package larkcli

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/zhangga/aino/pkg/logger"
)

type Client struct {
	ctx        context.Context
	token      *accessToken
	httpClient *http.Client
	larkClient *lark.Client
}

func NewClient(ctx context.Context, appId, appSecret string) *Client {
	return &Client{
		ctx:        ctx,
		token:      newAccessToken(appId, appSecret),
		httpClient: &http.Client{Timeout: 30 * time.Second},
		larkClient: lark.NewClient(appId, appSecret),
	}
}

func (c *Client) GetMessage(msgId string) (*larkim.GetMessageRespData, error) {
	// 创建请求对象
	req := larkim.NewGetMessageReqBuilder().
		MessageId(msgId).
		Build()
	// 发起请求
	resp, err := c.larkClient.Im.V1.Message.Get(c.ctx, req)
	// 处理错误
	if err != nil {
		logger.Errorf("请求失败: %v", err)
		return nil, err
	}
	// 服务端错误处理
	if !resp.Success() {
		logger.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		return nil, fmt.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}
	var msg LarkMessage
	if err = sonic.Unmarshal(resp.RawBody, &msg); err != nil {
		logger.Errorf("解析响应失败: %v", err)
		return nil, err
	}
	return resp.Data, nil
}

func (c *Client) GetMessageV1(msgId string) (*LarkMessage, error) {
	url := fmt.Sprintf("https://open.larkoffice.com/open-apis/im/v1/messages/%s", msgId)
	body, err := c.doRequest("GET", url, nil, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	})
	if err != nil {
		logger.Errorf("请求失败: %v", err)
		return nil, err
	}
	var msg LarkMessage
	if err = sonic.Unmarshal(body, &msg); err != nil {
		logger.Errorf("解析响应失败: %v", err)
		return nil, err
	}
	return &msg, nil
}

func (c *Client) SendMessage(receiveId, msgType, content, uuid string) error {
	return nil
}

func (c *Client) doRequest(method, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	token, err := c.token.ensureToken()
	if err != nil {
		logger.Errorf("获取access token失败: %v", err)
		return nil, err
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API返回错误状态码: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
