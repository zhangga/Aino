package conf

import logger "github.com/zhangga/aino/pkg/zlog"

var GlobalConfig Config

type IConfig interface {
	ValidData() bool
}
type Config struct {
	LogConf     *logger.Config `mapstructure:"log" json:"log" yaml:"log"`       // 日志相关配置
	EmbedConfig *EmbedConfig   `mapstructure:"embed" json:"embed" yaml:"embed"` // 向量配置
}

type EmbedConfig struct {
	BaseURL string `mapstructure:"base_url" json:"base_url" yaml:"base_url"` // 向量服务地址
	APIKey  string `mapstructure:"api_key" json:"api_key" yaml:"api_key"`    // 向量服务API Key
	Model   string `mapstructure:"model" json:"model" yaml:"model"`          // 向量模型
}

func (c *Config) ValidData() bool {
	return true
}
