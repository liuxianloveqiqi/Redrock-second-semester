package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr:     "43.139.195.17:6333",
		Password: "",
		DB:       0,
	})

	// 订阅channel1这个channel
	sub := rdb.Subscribe(ctx, "channel_1")
	// sub.Channel() 返回go channel，可以循环读取redis服务器发过来的消息
	for msg := range sub.Channel() {
		// 打印收到的消息
		fmt.Println(msg.Channel)
		fmt.Println(msg.Payload)
	}
}
