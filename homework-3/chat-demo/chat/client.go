// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chat

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512

	//心跳周期
	heartBeatPeriod = 10 * time.Second

	// 最大等待时间
	maxWaitTime = 60 * time.Second
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id  int64
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// 最后发送消息的时间
	lastMessageTime int64

	// 互斥锁
	mutex sync.Mutex

	// 限制发送消息的计时器
	limitSpeak *time.Ticker

	//敏感词列表
	sensitiveWords []string
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	c.hub.heartBeat = time.NewTicker(heartBeatPeriod)
	c.limitSpeak = time.NewTicker(3 * time.Second)
	lastMessageTime := time.Now()

	for {

		select {

		// 每10秒触发
		case <-c.hub.heartBeat.C:
			// 检测心跳，执行重连
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				log.Printf("掉线了，断开连接%v", err)
				c.reconnect()
				return
			}
		// 限制三秒才能发一次言
		case <-c.limitSpeak.C:
			// 每三秒重置最后发送消息的时间
			fmt.Println("请等待")
			lastMessageTime = time.Now()
		default:
			// 检查距离上一次发送消息的时间是否超过了3秒
			if time.Now().Sub(lastMessageTime) < 3*time.Second {
				// 超过了3秒才允许发送
				break
			}
			// 更新最后一次消息时间
			lastMessageTime = time.Now()
			// 处理图片数据
			_, message, err := c.conn.ReadMessage()
			if strings.HasPrefix(string(message), "image/") {
				c.hub.broadcast <- message
			} else {

				if err != nil {
					if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						log.Printf("error: %v", err)
					}
					break
				}
				message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
				// 敏感词替换
				c.sensitiveWords = []string{"傻逼死了"}
				if c.containsSensitiveWords(message) {
					message = c.replaceSensitiveWords(message)
				}
			}
			c.hub.broadcast <- message
			// 更新最后一次消息时间,使用原子操作更新最后发送消息时间防止并发出错
			atomic.StoreInt64(&c.lastMessageTime, time.Now().Unix())
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			// 对 send 进行写入时加锁
			c.mutex.Lock()
			w.Write(message)
			c.mutex.Unlock()
			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// 掉线重连
func (c *Client) reconnect() {
	for {
		time.Sleep(3 * time.Second)
		c.mutex.Lock() // 加锁
		waitTime := time.Now().Unix() - c.lastMessageTime
		c.mutex.Unlock() // 解锁

		if waitTime > int64(maxWaitTime.Seconds()) {
			log.Println("开始重新连接")
			u := url.URL{Scheme: "ws", Host: *flag.String("addr", ":8080", "http service address"), Path: "/ws"}
			// 开始连接
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			c.conn = conn
			c.send = make(chan []byte, 256)
			log.Println("重连成功")
			go c.writePump()
			go c.readPump()
			// 成功跳调出循环
			break

		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

// 检测消息是否含有敏感词
func (c *Client) containsSensitiveWords(message []byte) bool {
	if len(c.sensitiveWords) == 0 {
		return false
	}
	for _, word := range c.sensitiveWords {
		if strings.Contains(string(message), word) {
			return true
		}
	}
	return false
}

// 替换消息中的敏感词为*
func (c *Client) replaceSensitiveWords(message []byte) []byte {
	for _, word := range c.sensitiveWords {
		message = bytes.Replace(message, []byte(word), []byte(strings.Repeat("*", len(word))), -1)
	}
	return message
}
