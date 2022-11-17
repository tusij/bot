package openapi

import (
	"context"
	"github.com/tusij/bot.git/modle/dto"
	"github.com/tusij/bot.git/token"
	"time"
)

type OpenAPI interface {
	Base
	WebsocketAPI
	MessageAPI
}
type Base interface {
	Setup(token *token.Token, timeoutMS, idleConns int) OpenAPI

	WithTimeout(duration time.Duration) OpenAPI
}

type WebsocketAPI interface {
	GetWSInfo(ctx context.Context, params map[string]string, body string) (*dto.WebsocketAP, error)
}

type MessageAPI interface {
	PostMessage(ctx context.Context, channelID string, msg *dto.MessageToCreate) (*dto.Message, error)
}
