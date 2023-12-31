package process

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/openapi"
	"log"
	"net/url"
	"qqBot/database"
	"qqBot/pkg/service"
	"strings"
	"time"
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

// ProcessATMessage is a function to process message
func (p Processor) ProcessATMessage(input string, data *dto.WSATMessageData) error {
	ctx := context.Background()
	cmd := message.ParseCommand(input)
	channelId := data.ChannelID

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
		p.sendReply(ctx, channelId, toCreate)
	case "21点":
		input = "你作为庄家，我作为闲家，我们一起玩21点吧"
		toCreate.Content, _ = service.G.GeneratorGPTContent(ctx, []string{input})
		if toCreate.Content == "" {
			log.Println("GeneratorGPTContent error")
			return nil
		}
		//插入消息进入21点环节
		p.sendReply(ctx, channelId, toCreate)
		err := p.saveMessage(ctx, data.Author.ID, input)
		if err != nil {
			return err
		}
	case "播放":
		log.Println("播放", input)
		input = cmd.Content
		_, err := p.Api.PostAudio(ctx, channelId, &dto.AudioControl{
			URL:    "https://soundcloud.com/essenger/essenger-and-cryoshell-as?utm_source=clipboard&utm_medium=text&utm_campaign=social_sharing",
			Status: 0,
			Text:   "As Above, So Below",
		})
		if err != nil {
			log.Println("PostAudio error", err)
			return err
		}
	case "天气":
		log.Println(cmd.Content)
		//log.Println("天气", cmd.Content)
		escaped := url.QueryEscape(cmd.Content)
		location, err := service.W.FetchLocation(escaped)
		if err != nil {
			log.Println("FetchLocation error", err)
			return err
		}
		weather, err := service.W.FetchWeather(location)
		if err != nil {
			log.Println("FetchLocation error", err)
			return err
		}
		weatherByteArray, err := json.Marshal(weather.Daily)
		if err != nil {
			log.Println("FetchLocation error", err)
			return err
		}
		weatherString := string(weatherByteArray)
		toCreate.Content = weatherString
		p.sendReply(ctx, data.ChannelID, toCreate)
	default:
		fmt.Println([]string{input}[0])
		toCreate.Content, _ = service.G.GeneratorGPTContent(ctx, []string{input})
		if toCreate.Content == "" {
			log.Println("GeneratorGPTContent error")
			return nil
		}
		p.sendReply(ctx, data.ChannelID, toCreate)
	}

	return nil
}

func (p Processor) ProcessMessage(input string, data *dto.WSMessageData) error {
	ctx := context.Background()
	userId := data.Author.ID
	exist, _ := p.checkExist(ctx, userId)
	log.Println("input", input)

	//如果存在相应key，说明在游戏环节之中
	toCreate := &dto.MessageToCreate{
		MsgID: data.ID,
	}
	//如果不存在相应key，说明不在游戏环节之中
	if !exist {
		log.Println("key not exist")
		toCreate.Content = "不在任何一个聊天环节中"
		p.sendReply(ctx, data.ChannelID, toCreate)
		return nil
	}
	if strings.Index(input, "退出") != -1 {
		log.Println("退出")
		err := p.deleteMessageHistory(ctx, userId)
		if err != nil {
			return err
		}
		toCreate.Content = "退出之前的聊天环节"
		p.sendReply(ctx, data.ChannelID, toCreate)
		return nil
	}
	messages, _ := p.getMessageHistory(ctx, userId)
	log.Println("messages", messages)
	messages = append(messages, input)
	toCreate.Content, _ = service.G.GeneratorGPTContent(ctx, messages)
	p.sendReply(ctx, data.ChannelID, toCreate)
	err := p.saveMessage(ctx, userId, input)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (p Processor) sendReply(ctx context.Context, channelID string, toCreate *dto.MessageToCreate) {
	log.Println(toCreate)
	if _, err := p.Api.PostMessage(ctx, channelID, toCreate); err != nil {
		log.Println(err)
	}
}

func (p *Processor) defaultReplyContent() (string, error) {
	return "你好我是QQ机器人" + message.Emoji(307), nil
}

// 以用户的id为纬度，保存用户的消息，设置相应的时间
func (p *Processor) saveMessage(ctx context.Context, userID string, message string) error {
	// Redis key to store the list of user inputs
	key := database.GenerateKey(userID)
	expiration := 10 * time.Minute // 设置过期时间为10分钟
	err := database.RedisEngine.LPush(ctx, key, message).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	//设置过期时间
	err = database.RedisEngine.Expire(ctx, key, expiration).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// 以用户的id为纬度，获取用户的消息,在一定的时间范围内
func (p *Processor) getMessageHistory(ctx context.Context, userID string) ([]string, error) {
	// Redis key to store the list of user inputs
	key := database.GenerateKey(userID)
	values, err := database.RedisEngine.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return values, nil
}

// 以用户的id为纬度，删除用户的消息
func (p *Processor) deleteMessageHistory(ctx context.Context, userID string) error {
	// Redis key to store the list of user inputs
	key := database.GenerateKey(userID)
	err := database.RedisEngine.Del(ctx, key).Err()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (p *Processor) checkExist(ctx context.Context, userID string) (bool, error) {
	// Redis key to store the list of user inputs
	key := database.GenerateKey(userID)
	exist, err := database.RedisEngine.Exists(ctx, key).Result()
	if err != nil {
		log.Println(err)
		return false, err
	}
	return exist == 1, nil
}
