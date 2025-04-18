package conf

import "github.com/zhangga/aino/pkg/config"

var _ config.IConfig = (*Config)(nil)

type Config struct {
	ServiceConfig *ServiceConfig `yaml:"service" json:"service" mapstructure:"service"`
	// Lark配置
	LarkConfig *LarkConfig `yaml:"lark" json:"lark" mapstructure:"lark"`
	// LLM配置
	LLMConfig *LLMConfig `yaml:"llm" json:"llm" mapstructure:"llm"`
}

// ValidData implements config.IConfig.
func (c *Config) ValidData() bool {
	return true
}

type ServiceConfig struct {
	HttpPort int `yaml:"http_port" json:"http_port" mapstructure:"http_port"`
}

type LarkConfig struct {
	AppID     string `yaml:"app_id" json:"app_id" mapstructure:"app_id"`
	AppSecret string `yaml:"app_secret" json:"app_secret" mapstructure:"app_secret"`
}

type LLMConfig struct {
	Model   string `yaml:"model" json:"model" mapstructure:"model"`
	ApiKey  string `yaml:"api_key" json:"api_key" mapstructure:"api_key"`
	BaseURL string `yaml:"base_url" json:"base_url" mapstructure:"base_url"`
}
