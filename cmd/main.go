package main

import (
	"log"
	"qqBot/global"
	"qqBot/pkg/config"
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
