package main

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"math/rand"
	"sync"
	"time"
)

const (
	topic         = "PeakShaving_topic"
	channel       = "PeakShaving_channel"
	consumerCount = 5
	maxInFlight   = 10
)

type myHandler struct{}

func (*myHandler) HandleMessage(msg *nsq.Message) error {
	fmt.Printf("接收消息: %v\n", string(msg.Body))
	// 模拟消费者在处理消息时的随机延迟
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return nil
}

func main() {
	// 创建一个生产者
	producer, err := nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig())
	if err != nil {
		fmt.Printf("创建生产者失败: %v\n", err)
	}

	// 启动一个消费者组
	wg := sync.WaitGroup{}
	wg.Add(consumerCount)
	// 启动多个消费者
	var consumers []*nsq.Consumer
	for i := 0; i < consumerCount; i++ {
		consumer, err := nsq.NewConsumer(topic, channel, nsq.NewConfig())
		if err != nil {
			fmt.Printf("创建消费者失败: %v\n", err)
		}
		// 注册消费者的处理函数
		consumer.AddHandler(&myHandler{})
		// 连接到NSQD节点
		if err := consumer.ConnectToNSQD("127.0.0.1:4150"); err != nil {
			fmt.Printf("连接NSQ失败: %v\n", err)
		}
		consumers = append(consumers, consumer)
	}
	// 消费者并发消费
	for _, consumer := range consumers {
		go func(consumer *nsq.Consumer) {
			defer wg.Done()
			<-consumer.StopChan
			fmt.Println("消费者停止")
		}(consumer)
	}

	// 发送一些消息到NSQ中
	for i := 0; i < 100; i++ {
		msg := fmt.Sprintf("message %d\n", i)
		if err := producer.Publish(topic, []byte(msg)); err != nil {
			fmt.Printf("发送消息失败: %v\n", err)
			time.Sleep(time.Second)
			continue
		}
		fmt.Printf("发送消息: %v\n", msg)
		// 模拟生产者在处理消息时的随机延迟
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	}

	// 停止生产者
	producer.Stop()
	// 停止消费者
	for _, consumer := range consumers {
		consumer.Stop()
	}
	wg.Wait()
	fmt.Println("程序退出")
}
