package main

import (
	"ResourceKeeper/ConfigManage"
	"fmt"
	RPC "github.com/rpcclient/RPCClient_Init/Client"
	"github.com/rpcclient/RPCClient_Init/Const"
)

func main() {
	fmt.Println("------------------------------------------------------------------------Webconfig初始化------------------------------------------------------------------------")
	err := ConfigManage.ConfigInit("ConfigManage/config.yaml") // 这里假设配置文件名为 config.yaml
	if err != nil {
		fmt.Println("Failed to load config: %v\n", err)
		fmt.Printf("Failed to load config: %v\n", err)
		return
	}

	fmt.Println("------------------------------------------------------------------------RPC初始化------------------------------------------------------------------------")
	RPC.RpcClientMapInit()
	go RPC.Init_Client.RpcClient_Init("127.0.0.1:8080", Const.VRTS_SERVER_TYPE_WEBSERVER)
}
