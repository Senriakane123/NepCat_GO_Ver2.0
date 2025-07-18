package main

import (
	"RPCServer/RpcServer_Init/Server"
	"sync"
)

var wg sync.WaitGroup

func main() {
	serverapi := Server.RpcServerMapInit()
	go serverapi.RpcServer_Init("127.0.0.1:8080")
	wg.Add(1)
	wg.Wait()
}
