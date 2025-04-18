package rolelist

import (
	"context"
	"github.com/spf13/viper"
	"github.com/zhangga/aino/pkg/logger"
)

const roleListFilePath = "configs/rolelist.yaml"

type Role struct {
	Title   string   `yaml:"title" mapstructure:"title"`
	Content string   `yaml:"content" mapstructure:"content"`
	Tags    []string `yaml:"tags" mapstructure:"tags"`
}

type RoleList struct {
	RoleList []*Role `yaml:"rolelist" mapstructure:"rolelist"`
}

var roleList RoleList

func InitRoleList(ctx context.Context) {
	roleViper := viper.New()
	roleViper.SetConfigFile(roleListFilePath)
	if err := roleViper.ReadInConfig(); err != nil {
		logger.Fatalf("read rolelist file path=%s, error: %v", roleListFilePath, err)
	}
	if err := roleViper.Unmarshal(&roleList); err != nil {
		logger.Fatalf("unmarshal rolelist file path=%s, error: %v", roleListFilePath, err)
	}
}
