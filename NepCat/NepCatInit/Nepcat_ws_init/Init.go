package Nepcat_ws_init

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type NepcatWebSocket struct {
	conn           *websocket.Conn
	messageChannel chan string
}

var NepcatWS *NepcatWebSocket

func WebChannelInit() {
	NepcatWS = &NepcatWebSocket{
		messageChannel: make(chan string, 100),
	}
}

// WebSocketInit 会持续尝试连接直到成功
func WebSocketInit(Scheme, host string, port int, path, raw string) {
	for {
		fmt.Println("🔄 尝试连接 WebSocket...")

		// 构建 WebSocket URL
		serverURL := url.URL{
			Scheme:   Scheme,
			Host:     host + ":" + strconv.Itoa(port),
			Path:     path,
			RawQuery: raw,
		}

		// 建立连接
		conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
		if err != nil {
			log.Printf("❌ 连接失败: %v，5 秒后重试...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// 成功建立连接
		fmt.Println("✅ 成功连接到 WebSocket 服务器:", serverURL.String())

		// 初始化全局实例
		NepcatWS.conn = conn

		// 捕获 Ctrl+C 退出
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt)

		// 启动接收协程
		go func() {
			for {
				_, message, err := NepcatWS.conn.ReadMessage()
				if err != nil {
					log.Println("❌ WebSocket 连接中断，尝试重连:", err)
					break // 跳出接收循环，重新进入连接流程
				}
				NepcatWS.messageChannel <- string(message)
			}
		}()

		// 主线程等待中断信号
		<-interrupt
		log.Println("⏳ 收到中断信号，准备关闭 WebSocket")
		NepcatWS.conn.Close()
		return
	}
}

// 获取消息通道
func (ws *NepcatWebSocket) GetChannel() *chan string {
	return &ws.messageChannel
}

func (ws *NepcatWebSocket) GetConn() *websocket.Conn {
	return ws.conn
}
