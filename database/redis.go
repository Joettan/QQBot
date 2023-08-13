package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"qqBot/global"
)

func InitRedisEngine(ctx context.Context) {
	host := global.RedisConfig.Host
	port := global.RedisConfig.Port
	password := global.RedisConfig.Password
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port, // Redis 地址
		Password: password,          // no password set
		DB:       0,                 // use default DB
	})
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	log.Println(pong)

	RedisEngine = rdb
}

func GenerateKey(userID string) string {
	return "user:inputs" + userID
}
