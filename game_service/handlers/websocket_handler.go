package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 用于存储所有活动连接的全局变量
var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 默认接受所有来源的连接
		return true
	},
}

// WebSocket 连接处理
func WebSocketHandler(c *gin.Context) {
	// 升级 HTTP 请求为 WebSocket 连接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error: ", err)
		return
	}
	defer conn.Close()

	// 将连接添加到客户端列表
	clients[conn] = true
	fmt.Println("New WebSocket connection established")

	// 处理来自客户端的消息
	for {
		// 等待客户端发送消息
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// 如果连接关闭或发生错误，退出循环
			fmt.Println("Connection closed or error: ", err)
			delete(clients, conn)
			break
		}

		// 如果接收到的是文本消息
		if messageType == websocket.TextMessage {
			// 先判断消息是否为普通字符串
			if isValidString(p) {
				// 如果是普通字符串，直接处理
				fmt.Println("Received plain string message:", string(p))
				response := Message{
					Type:    "response",
					Content: "Received plain string: " + string(p),
				}
				responseJSON, _ := json.Marshal(response)
				conn.WriteMessage(websocket.TextMessage, responseJSON)
			} else {
				// 如果不是普通字符串，尝试解析为 JSON
				var msg map[string]interface{}
				err := json.Unmarshal(p, &msg)
				if err != nil {
					// 如果解析失败，返回错误信息
					fmt.Println("Error parsing message as JSON:", err)
					response := Message{
						Type:    "error",
						Content: "Invalid message format",
					}
					responseJSON, _ := json.Marshal(response)
					conn.WriteMessage(websocket.TextMessage, responseJSON)
				} else {
					// 解析成功，处理 JSON 消息
					fmt.Println("Received JSON message:", msg)
					response := Message{
						Type:    "response",
						Content: "Received JSON message",
					}
					responseJSON, _ := json.Marshal(response)
					conn.WriteMessage(websocket.TextMessage, responseJSON)
				}
			}
		}
	}
}

// 判断是否是有效的普通字符串
func isValidString(p []byte) bool {
	// 判断是否为普通字符串的简单规则，可以根据需求更改
	// 这里仅简单判断是否不是空字符串
	return len(p) > 0 && p[0] != '{' && p[0] != '[' // 不是 JSON 字符串
}

// Message 是发送到客户端的 JSON 数据结构
type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}
