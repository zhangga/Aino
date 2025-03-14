package larkcli

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/bytedance/sonic"
)

const tokenExpiry = 5 * time.Minute

type accessToken struct {
	sync.RWMutex
	appId     string
	appSecret string
	token     string
	expiry    time.Time
}

func newAccessToken(appId, appSecret string) *accessToken {
	return &accessToken{
		appId:     appId,
		appSecret: appSecret,
		expiry:    time.Now(),
	}
}

func (t *accessToken) ensureToken() (string, error) {
	t.RLock()
	// 如果token未过期，直接返回
	if time.Until(t.expiry) > tokenExpiry {
		defer t.RUnlock()
		return t.token, nil
	}

	// 更新token
	t.RUnlock()
	t.Lock()
	defer t.Unlock()

	// 再次检查token是否过期，防止在获取锁和更新token之间有其他goroutine更新了token
	if time.Until(t.expiry) > tokenExpiry {
		return t.token, nil
	}

	// 这里添加实际的token更新逻辑
	req, err := http.NewRequest("POST",
		"https://open.larkoffice.com/open-apis/auth/v3/tenant_access_token/internal",
		strings.NewReader(url.Values{"app_id": {t.appId}, "app_secret": {t.appSecret}}.Encode()))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求执行失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("获取access token失败 resp.StatusCode: %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)

	var tokenResp LarkAccessTokenResp
	if err := sonic.Unmarshal(respBody, &tokenResp); err != nil {
		return "", fmt.Errorf("解析token响应失败(原始响应:%.200q): %w", string(respBody), err)
	}
	if tokenResp.Code != 0 || tokenResp.Expire == 0 || tokenResp.Token == "" {
		return "", fmt.Errorf("飞书API返回无效数据结构(响应码:%d 内容:%.200q)", resp.StatusCode, string(respBody))
	}
	t.token = tokenResp.Token
	t.expiry = time.Now().Add(time.Duration(tokenResp.Expire) * time.Second)
	return t.token, nil
}

type LarkAccessTokenResp struct {
	Code   int    `json:"code"`
	Expire int    `json:"expire"`
	Msg    string `json:"msg"`
	Token  string `json:"tenant_access_token"`
}
