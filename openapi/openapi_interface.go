package openapi

import (
	"context"
	"github.com/tusij/bot.git/config"
	"github.com/tusij/bot.git/modle/dto"
	"time"
)

type OpenAPI interface {
	Base
	WebsocketAPI
	MessageAPI
}
type Base interface {
	Setup(c *config.OpenAPIConfig) OpenAPI

	WithTimeout(duration time.Duration) OpenAPI
}

type WebsocketAPI interface {
	GetWSInfo(ctx context.Context, params map[string]string, body string) (*dto.WebsocketAP, error)
}

type MessageAPI interface {
	PostMessage(ctx context.Context, channelID string, msg *dto.MessageToCreate) (*dto.Message, error)
}
