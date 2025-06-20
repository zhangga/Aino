package handlers

import (
	"context"
)

func NewClearCardHandler() (CardKind, CardHandlerFunc) {
	return ClearCardKind, func(ctx context.Context, m *MessageHandler, info *CardActionInfo) error {
		return nil
	}
}
