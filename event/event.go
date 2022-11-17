package event

import (
	"encoding/json"
	"errors"
	"github.com/tidwall/gjson" // message中含有d字段 d的结构不定 强行解析会有问题
	"github.com/tusij/bot.git/modle/dto"
)

type eventParseFunc func(event *dto.WSPayload, message []byte) error

var eventParseFuncMap = map[dto.OPCode]map[dto.EventType]eventParseFunc{
	dto.WSDispatchEvent: {
		dto.EventAtMessageCreate: atMessageHandler,
	},
}

func ParseData(message []byte, target interface{}) error {
	data := gjson.Get(string(message), "d")
	return json.Unmarshal([]byte(data.String()), target)
}

func ParseAndHandle(payload *dto.WSPayload) error {
	// 指定类型的 handler
	if h, ok := eventParseFuncMap[payload.OPCode][payload.Type]; ok {
		return h(payload, payload.RawMessage)
	}

	return errors.New("not found func")
}

func atMessageHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.Message{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if DefaultHandlers.ATMessage != nil {
		return DefaultHandlers.ATMessage(payload, data)
	}
	return errors.New("not found func")
}
