package handlers

import (
	"context"
)

type CardHandlerMeta func() (CardKind, CardHandlerFunc)

type CardHandlerFunc func(ctx context.Context, m *MessageHandler, info *CardActionInfo) error
