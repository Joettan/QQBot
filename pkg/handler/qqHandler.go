package handler

import (
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"qqBot/pkg/process"
	"strings"
)

func MessageHandler() event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		//fmt.Println(event, data)
		input := strings.ToLower(message.ETLInput(data.Content))
		return process.P.ProcessMessage(input, data)
	}
}
