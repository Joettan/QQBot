package handler

import (
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"log"
	"qqBot/pkg/process"
	"strings"
)

func ATMessageHandler() event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		//fmt.Println(event, data)
		input := strings.ToLower(message.ETLInput(data.Content))
		return process.P.ProcessATMessage(input, data)
	}
}

func MessageHandler() event.MessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content))
		log.Println(input)
		return process.P.ProcessMessage(input, data)
	}
}
