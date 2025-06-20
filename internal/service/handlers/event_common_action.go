package handlers

import (
	"errors"
	"github.com/zhangga/aino/internal/rolelist"
	"github.com/zhangga/aino/internal/utils"
)

var (
	_ Action = (*ProcessedUniqueAction)(nil)
	_ Action = (*ClearAction)(nil)
	_ Action = (*RoleListAction)(nil)
)

var (
	ErrProcessedUniqueAction = errors.New("processed unique action")
	ErrClearAction           = errors.New("clear action")
	ErrRoleListAction        = errors.New("role list action")
)

// ProcessedUniqueAction 避免重复处理
type ProcessedUniqueAction struct {
}

func (ProcessedUniqueAction) Execute(data *ActionData) error {
	if data.handler.msgCache.IfProcessed(data.info.GetMsgId()) {
		return ErrProcessedUniqueAction
	}
	data.handler.msgCache.TagProcessed(data.info.GetMsgId())
	return nil
}

type ClearAction struct {
}

func (ClearAction) Execute(a *ActionData) error {
	if _, found := utils.EitherTrimEqual(a.info.GetContent(), "/clear", "清除"); found {
		if err := sendClearCacheCheckCard(a.ctx, a.info.GetSessionId(), a.info.GetMsgId()); err != nil {
			return err
		}
		return ErrClearAction
	}
	return nil
}

type RoleListAction struct {
}

func (r RoleListAction) Execute(a *ActionData) error {
	if _, foundSystem := utils.EitherTrimEqual(a.info.GetContent(), "/roles", "角色列表"); foundSystem {
		tags := rolelist.GetAllUniqueTags()
		SendRoleTagsCard(a.ctx, a.info.GetSessionId(), a.info.GetMsgId(), tags)
		return ErrRoleListAction
	}
	return nil
}
