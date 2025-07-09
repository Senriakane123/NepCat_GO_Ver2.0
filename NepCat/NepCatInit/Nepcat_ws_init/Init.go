package Nepcat_ws_init

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
)

type NepcatWebSocket struct {
	conn           *websocket.Conn
	messageChannel chan string
}

var NepcatWS *NepcatWebSocket

//func Init() {
//	if NepcatWS == nil {
//		NepcatWS.WebSocketInit()
//	}
//}

// WebSocketInit 函数用于初始化 WebSocket 连接
func WebSocketInit(Scheme, host string, port int, path, raw string) {
	NepcatWS.messageChannel = make(chan string, 100)

	// 创建 WebSocket 服务器 URL
	serverURL := url.URL{
		Scheme:   Scheme,
		Host:     host + string(port),
		Path:     path,
		RawQuery: raw,
	}

	var err error
	// 使用默认的 WebSocket Dialer 连接服务器
	NepcatWS.conn, _, err = websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		// 如果连接失败，打印错误信息并退出程序
		log.Fatalf("❌ 连接 WebSocket 失败: %v", err)
	}
	fmt.Println("✅ 成功连接到 WebSocket 服务器")

	// 捕获 Ctrl+C 退出
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// 启动一个 goroutine，用于读取服务器发送的消息
	go func() {
		for {
			_, message, err := NepcatWS.conn.ReadMessage()
			if err != nil {
				// 如果读取消息失败，打印错误信息并退出 goroutine
				log.Println("❌ 读取消息失败:", err)
				return
			}
			// 打印收到的消息
			fmt.Println("📩 收到消息:", string(message))
			// 将收到的消息发送到消息通道
			NepcatWS.messageChannel <- string(message)
			//go MessageHandler()
		}
	}()

	// 等待 Ctrl+C 退出
	<-interrupt
	fmt.Println("⏳ 关闭 WebSocket 连接...")
}

// GetChannel 返回消息通道
func (ws *NepcatWebSocket) GetChannel() *chan string {
	return &ws.messageChannel
}

func (ws *NepcatWebSocket) GetConn() *websocket.Conn {
	return ws.conn
}
