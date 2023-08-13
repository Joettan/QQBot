package main

import (
	"context"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"log"
	"qqBot/global"
	"qqBot/pkg/config"
	"qqBot/pkg/handler"
	"qqBot/pkg/process"
	"qqBot/pkg/service"
	"time"
)

func main() {

	token := token.BotToken(global.BotConfig.AppID, global.BotConfig.Token)
	api := botgo.NewOpenAPI(token).WithTimeout(3 * time.Second)
	ctx := context.Background()

	ws, err := api.WS(ctx, nil, "")
	log.Printf("%+v, err:%v", ws, err)

	process.InitProcessor(api)

	//注册服务
	_ = service.NewFactory()

	intent := websocket.RegisterHandlers(handler.MessageHandler())

	log.Printf("intent:%+v", intent)

	// 启动 session manager 进行 ws 连接的管理，如果接口返回需要启动多个 shard 的连接，这里也会自动启动多个
	botgo.NewSessionManager().Start(ws, token, &intent)
}

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err :%v", err)
		panic(err)
	}
	service.InitGPTService()
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
	err = setting.ReadSection("GPTConfig", &global.GPTConfig)
	if err != nil {
		return err
	}
	return nil
}
