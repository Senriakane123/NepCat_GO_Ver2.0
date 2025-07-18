package main

import (
	"RpcClient/RPCClient_Init/Client"
	"RpcClient/RPCClient_Init/Const"
	"sync"
)

var wg sync.WaitGroup

func main() {
	Client.RpcClientMapInit()
	go Client.Init_Client.RpcClient_Init("127.0.0.1:8080", Const.VRTS_SERVER_TYPE_WEBSERVER)
	wg.Add(1)
	wg.Wait()
}
