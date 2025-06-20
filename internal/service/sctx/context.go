package sctx

import (
	"context"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/zhangga/aino/internal/conf"
)

type serviceConfigKeyType struct{}
type larkClientKeyType struct{}

var (
	serviceConfigKey = serviceConfigKeyType{}
	larkClientKey    = larkClientKeyType{}
)

type Context struct {
	context.Context
}

func WithContext(ctx context.Context) *Context {
	return &Context{Context: ctx}
}

func (c *Context) WithServiceConfig(config *conf.ServiceConfig) *Context {
	return WithContext(context.WithValue(c, serviceConfigKey, config))
}

func GetServiceConfig(ctx context.Context) *conf.ServiceConfig {
	return ctx.Value(serviceConfigKey).(*conf.ServiceConfig)
}

func (c *Context) WithLarkClient(client *lark.Client) *Context {
	return WithContext(context.WithValue(c, larkClientKey, client))
}

func GetLarkClient(ctx context.Context) *lark.Client {
	return ctx.Value(larkClientKey).(*lark.Client)
}
