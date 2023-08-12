package main

import (
	"context"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"log"
	"qqBot/global"
	"qqBot/pkg/config"
	"qqBot/pkg/process"
	"strings"
	"time"
)

var (
	p process.Processor
)

func main() {

	//client := openai.NewClient("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")
	//
	//resp, err := client.CreateChatCompletion(
	//	context.Background(),
	//	openai.ChatCompletionRequest{
	//		Model: openai.GPT3Dot5Turbo,
	//		Messages: []openai.ChatCompletionMessage{
	//			{
	//				Role:    openai.ChatMessageRoleUser,
	//				Content: "Hello!",
	//			},
	//		},
	//	},
	//)
	//
	//if err != nil {
	//	fmt.Printf("ChatCompletion error: %v\n", err)
	//	return
	//}
	//
	//fmt.Println(resp.Choices[0].Message.Content)
	//fmt.Println("test :", global.BotConfig.AppID)
	token := token.BotToken(global.BotConfig.AppID, global.BotConfig.Token)
	api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()

	ws, err := api.WS(ctx, nil, "")
	log.Printf("%+v, err:%v", ws, err)

	p = process.Processor{Api: api}
	//fmt.Println("test")

	// 监听哪类事件就需要实现哪类的 handler，定义：websocket/event_handler.go
	//var atMessage event.ATMessageEventHandler = func(event *dto.WSPayload, data *dto.WSATMessageData) error {
	//	fmt.Println(event, data)
	//	return nil
	//}

	//var message event.MessageEventHandler = func(event *dto.WSPayload, data *dto.WSMessageData) error {
	//	fmt.Println(event, data)
	//	return nil
	//}

	intent := websocket.RegisterHandlers(MessageHandler())

	// 启动 session manager 进行 ws 连接的管理，如果接口返回需要启动多个 shard 的连接，这里也会自动启动多个
	botgo.NewSessionManager().Start(ws, token, &intent)
}

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err :%v", err)
		panic(err)
	}
}
func setupSetting() error {
	setting, err := config.NewSetting()
	if err != nil {
		return err
	}
	err = setting.ReadSection("QQBot", &global.BotConfig)
	if err != nil {
		return err
	}
	return nil
}

func MessageHandler() event.MessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageData) error {
		//fmt.Println(event, data)
		input := strings.ToLower(message.ETLInput(data.Content))
		return p.ProcessMessage(input, data)
	}
}
