package process

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
	"log"
)

type Processor struct {
	Api openapi.OpenAPI
}

var (
	P *Processor
)

func InitProcessor(api openapi.OpenAPI) error {
	P = &Processor{
		Api: api,
	}
	return nil
}

// ProcessMessage is a function to process message
func (p Processor) ProcessMessage(input string, data *dto.WSATMessageData) error {
	ctx := context.Background()
	cmd := message.ParseCommand(input)

	//区分的是相应内容，而不是消息的header
	toCreate := &dto.MessageToCreate{
		MsgID: data.ID,
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
	}

	switch cmd.Cmd {
	case "hi":
		toCreate.Content, _ = p.defaultReplyContent()
		p.sendReply(ctx, data.ChannelID, toCreate)
	default:
		toCreate.Content, _ = p.generatorGPTContent(ctx, input)
		p.sendReply(ctx, data.ChannelID, toCreate)
	}

	return nil
}

func (p Processor) sendReply(ctx context.Context, channelID string, toCreate *dto.MessageToCreate) {
	if _, err := p.Api.PostMessage(ctx, channelID, toCreate); err != nil {
		log.Println(err)
	}
}

func (p *Processor) defaultReplyContent() (string, error) {
	return "你好我是QQ机器人" + message.Emoji(307), nil
}

func (p *Processor) generatorGPTContent(ctx context.Context, msg string) (string, error) {
	client := openai.NewClient("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")
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)
	if err != nil {
		log.Printf("ChatGPT queryApi :%v", err)
	}

	return resp.Choices[0].Message.Content, nil
}
