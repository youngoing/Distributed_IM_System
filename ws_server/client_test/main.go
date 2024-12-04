package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
	"ws_server/shared"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// 发送消息的函数
func sendMessages(conn *websocket.Conn, userId string) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		var receiverIds []string
		var msgContent string
		var groupId string
		logrus.Println("确认群聊或私聊，输入1为群聊，输入2为私聊")
		if scanner.Scan() {
			chatType := scanner.Text()
			if chatType == "1" {
				logrus.Println("群聊")
				logrus.Print("请输入接收者userId（多个接收者用逗号分隔）: ")
				// 获取接收者 userId
				if scanner.Scan() {
					receiverIdsInput := scanner.Text()
					if receiverIdsInput == "" {
						logrus.Println("接收者userId不能为空，请重新输入")
						continue
					}
					receiverIds = strings.Split(receiverIdsInput, ",")
				} else {
					if err := scanner.Err(); err != nil {
						logrus.Println("读取接收者userId时出错:", err)
					}
					continue
				}

				logrus.Print("请输入groupId: ")
				// 获取 groupId
				if scanner.Scan() {
					groupId = scanner.Text()
					if groupId == "" {
						logrus.Println("groupId不能为空，请重新输入")
						continue
					}
				} else {
					if err := scanner.Err(); err != nil {
						logrus.Println("读取groupId时出错:", err)
					}
					continue
				}

			} else if chatType == "2" {
				logrus.Println("私聊")
				logrus.Print("请输入接收者userId: ")
				// 获取接收者 userId
				if scanner.Scan() {
					receiverId := scanner.Text()
					if receiverId == "" {
						logrus.Println("接收者userId不能为空，请重新输入")
						continue
					}
					receiverIds = append(receiverIds, receiverId)
				} else {
					if err := scanner.Err(); err != nil {
						logrus.Println("读取接收者userId时出错:", err)
					}
					continue
				}
			} else {
				logrus.Println("输入错误")
				continue
			}
		} else {
			if err := scanner.Err(); err != nil {
				logrus.Println("读取聊天类型时出错:", err)
			}
			continue
		}

		logrus.Print("请输入消息内容: ")
		if scanner.Scan() {
			msgContent = scanner.Text()
			if msgContent == "" {
				logrus.Println("消息内容不能为空，请重新输入")
				continue
			}
		} else {
			if err := scanner.Err(); err != nil {
				logrus.Println("读取消息内容时出错:", err)
			}
			continue
		}

		// 创建新的 WsMessage
		var newWsMessage shared.WsMsg
		if groupId != "" {
			newWsMessage = shared.NewWsGroupMessage(receiverIds, userId, groupId, msgContent)
		} else {
			newWsMessage = shared.NewWsUserMessage(userId, receiverIds, msgContent)
		}

		// 将消息转为 JSON 格式
		message, err := json.Marshal(newWsMessage)
		if err != nil {
			logrus.Println("消息转换为JSON格式失败:", err)
			continue
		}

		// 发送消息
		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			logrus.Println("发送消息失败:", err)
			return
		}
		logrus.Println("消息发送成功")
	}
}

func main() {
	userID := readUserId()
	logrus.Printf("用户ID: %s", userID)
	url := fmt.Sprint("ws://127.0.0.1:8001/?token=", userID)

	// 尝试连接 WebSocket 并自动重试
	var conn *websocket.Conn
	var err error
	count := 1

	for {
		if count > 3 {
			logrus.Println("连接到 WebSocket 失败，重试次数过多，程序退出")
			return
		}
		conn, _, err = websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			logrus.Printf("连接到 WebSocket 失败: %v", err)
			logrus.Printf("将在 %d 秒后重试...", count*2)          // 延迟的时间是递增的
			time.Sleep(time.Duration(count*2) * time.Second) // 重试等待时间
			count++                                          // 增加重试计数
		} else {
			logrus.Println("成功连接到 WebSocket")
			break // 连接成功，跳出重试循环
		}
	}
	defer conn.Close()

	// 启动接收消息的 goroutine
	go receiveMessages(conn)

	// 发送消息
	sendMessages(conn, userID)
}

// 接收消息的函数
func receiveMessages(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			logrus.Println("读取消息失败:", err)
			return
		}
		var msg shared.WsMsg
		err = json.Unmarshal(message, &msg)
		if err != nil {
			logrus.Println("解析消息失败:", err)
			continue
		}
		logrus.Println("收到消息:", msg.PrettyPrint())
	}
}

func readUserId() string {
	scanner := bufio.NewScanner(os.Stdin)
	logrus.Print("请输入用户ID: ")
	if scanner.Scan() {
		userId := scanner.Text()
		if userId == "" {
			logrus.Println("用户ID不能为空，请重新输入")
			return readUserId()
		}
		return userId
	}
	if err := scanner.Err(); err != nil {
		logrus.Println("读取用户ID时出错:", err)
	}
	return readUserId()
}
