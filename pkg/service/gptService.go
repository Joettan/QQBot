package service

import (
	"context"
	"fmt"
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

func (g *GptService) GeneratorGPTContent(ctx context.Context, msg []string) (string, error) {
	client := openai.NewClient(global.GPTConfig.GPTToken)
	messages := make([]openai.ChatCompletionMessage, 0, len(msg))
	for _, v := range msg {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: v,
		})
	}
	log.Println(messages)
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			//Messages: messages,
			Messages: messages,
		},
	)
	if err != nil {
		log.Printf("ChatGPT queryApi :%v", err)
		return "", err
	}
	fmt.Println(resp.Choices)

	return resp.Choices[0].Message.Content, nil
}
