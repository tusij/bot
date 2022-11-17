package event

import (
	"github.com/tusij/bot.git/modle/dto"
)

var DefaultHandlers struct {
	ATMessage ATMessageEventHandler
}

type ATMessageEventHandler func(event *dto.WSPayload, data *dto.Message) error

func RegisterHandlers(handlers ...interface{}) {
	for _, h := range handlers {
		switch handle := h.(type) {
		case ATMessageEventHandler:
			DefaultHandlers.ATMessage = handle
		default:
		}
	}
}
