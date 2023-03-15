package main

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "43.139.195.17:6333",
		Password: "",
		DB:       0,
	})
	// 将"hello World"消息发送到channel1这个通道上
	rdb.Publish(ctx, "channel_1", "hello World")
}
