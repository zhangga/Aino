package handlers

import (
	"context"
	"errors"
	"github.com/zhangga/aino/internal/rolelist"
	"github.com/zhangga/aino/pkg/logger"
)

var (
	ErrorCastCardActionInfo = errors.New("cast card action info failed")
	ErrorGetRole            = errors.New("get role failed")
)

func NewRoleTagCardHandler() (CardKind, CardHandlerFunc) {
	return RoleTagsChooseKind, func(ctx context.Context, m *MessageHandler, info *CardActionInfo) error {
		// 处理角色标签卡片
		return CommonProcessRoleTag(ctx, info)
	}
}

func CommonProcessRoleTag(ctx context.Context, info *CardActionInfo) error {
	logger.Debugf("role tags choose %v", info)
	tag := info.Option
	titles := rolelist.GetTitleListByTag(tag)
	SendRoleListCard(ctx, info.GetSessionId(), info.GetMsgId(), tag, titles)
	return nil
}

func NewRoleCardHandler() (CardKind, CardHandlerFunc) {
	return RoleChooseKind, func(ctx context.Context, m *MessageHandler, info *CardActionInfo) error {
		// 处理角色卡片
		return CommonProcessRoleCard(ctx, m, info)
	}
}

func CommonProcessRoleCard(ctx context.Context, m *MessageHandler, info *CardActionInfo) error {
	logger.Debugf("role cards choose %v", info)
	title := info.Option
	role := rolelist.GetFirstRoleByTitle(title)
	if role == nil {
		logger.Errorf("get role failed: %s", title)
		return ErrorGetRole
	}
	m.sessionCache.Clear(info.GetSessionId())

	return nil
}
