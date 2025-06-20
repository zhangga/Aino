package rolelist

import (
	"context"
	"github.com/duke-git/lancet/v2/slice"
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

func GetAllUniqueTags() []string {
	tags := make([]string, len(roleList.RoleList))
	for _, role := range roleList.RoleList {
		tags = append(tags, role.Tags...)
	}
	return slice.Union(tags)
}

func GetAllRoleTitles() []string {
	titles := make([]string, len(roleList.RoleList))
	for _, role := range roleList.RoleList {
		titles = append(titles, role.Title)
	}
	return titles
}

func GetFirstRoleByTitle(title string) *Role {
	for _, role := range roleList.RoleList {
		if role.Title == title {
			return role
		}
	}
	return nil
}

func GetFirstRoleByTag(tag string) *Role {
	for _, role := range roleList.RoleList {
		if slice.Contain(role.Tags, tag) {
			return role
		}
	}
	return nil
}

func GetTitleListByTag(tag string) []string {
	titles := make([]string, 0)
	for _, role := range roleList.RoleList {
		if slice.Contain(role.Tags, tag) {
			titles = append(titles, role.Title)
		}
	}
	return titles
}
