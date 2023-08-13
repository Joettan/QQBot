package service

import (
	"context"
	"github.com/sashabaranov/go-openai"
	"log"
	"qqBot/global"
)

var (
	G *GptService
)

type GptService struct {
	client *openai.Client
}

func InitGPTService() {
	G = &GptService{
		client: openai.NewClient(global.GPTConfig.GPTToken),
	}
}

func (g *GptService) generatorGPTContent(ctx context.Context, msg string) (string, error) {
	resp, err := g.client.CreateChatCompletion(
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
