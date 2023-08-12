package process

import (
	"context"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
	"log"
)

type Processor struct {
	Api openapi.OpenAPI
}

// ProcessMessage is a function to process message
func (p Processor) ProcessMessage(input string, data *dto.WSMessageData) error {
	ctx := context.Background()
	//cmd := message.ParseCommand(input)
	//toCreate := &dto.MessageToCreate{
	//	Content: "默认回复",
	//	//MessageReference: &dto.MessageReference{
	//	//	// 引用这条消息
	//	//	MessageID:             data.ID,
	//	//	IgnoreGetMessageError: true,
	//	//},
	//}

	// 进入到私信逻辑
	//if cmd.Cmd == "dm" {
	//	p.dmHandler(data)
	//	return nil
	//}

	//fmt.Println("Hi" + data.ChannelID)

	//switch cmd.Cmd {
	//case "hi":
	//	p.sendReply(ctx, data.ChannelID, toCreate)
	//default:
	//}

	//_, err := p.api.PostMessage(ctx, data.ChannelID, toCreate)
	//if err != nil {
	//	fmt.Println(err)
	//}

	_, _ = p.Api.PostMessage(ctx, data.ChannelID, &dto.MessageToCreate{
		MsgID:   data.ID, //如果未空则表示主动消息
		Content: "hello world",
	})

	return nil
}

func (p Processor) sendReply(ctx context.Context, channelID string, toCreate *dto.MessageToCreate) {
	if _, err := p.Api.PostMessage(ctx, channelID, toCreate); err != nil {
		log.Println(err)
	}
}
