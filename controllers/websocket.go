package controllers

import (
	"time"

	"github.com/gorilla/websocket"
)

type KenShinSocketController struct {
	BaseController
}

// socket
var upgrader = websocket.Upgrader{}

// 广播
var broadCast = make(chan Message)

// 链接的客户端
var clients = make(map[*websocket.Conn]bool)

type Message struct {
	Content string `json:"message"`
}

// 初始化并开启广播监听
func init() {
	go HandlerBroadCast()
}

func (ws *KenShinSocketController) Get() {

	conn, err := upgrader.Upgrade(ws.Ctx.ResponseWriter, ws.Ctx.Request, nil)

	if err != nil {
		return
	}

	defer conn.Close()

	// 添加到客户端集合内
	clients[conn] = true

	for {

		time.Sleep(3 * time.Second)

		SendMessage("发消息咯:" + time.Now().Format("2021-10-31 14:59:59"))

		// 获取消息
		var receivedMessage = Message{}

		jsErr := conn.ReadJSON(&receivedMessage)

		if jsErr != nil {
			println("断开链接!")
			closeClient(conn)
			break
		} else {

			println("received : ", receivedMessage.Content)
		}
	}
}

// 发送消息
func SendMessage(message string) {
	broadCast <- Message{Content: message}
}

// 广播信息处理
func HandlerBroadCast() {

	for {
		msg := <-broadCast
		for client := range clients {

			sendErr := client.WriteJSON(msg)

			if sendErr != nil {

				closeClient(client)
			}
		}
	}
}

// 关闭并且移除当前客户端
func closeClient(client *websocket.Conn) {

	println("关闭当前客户端")

	client.Close()

	delete(clients, client)
}
