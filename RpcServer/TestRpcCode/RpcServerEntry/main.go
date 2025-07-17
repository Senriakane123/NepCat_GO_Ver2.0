package main

import (
	"RPCServer/TestRpcCode/Server"
	"fmt"
	"log"
)

func main() {

	//var num int = 100
	//v := reflect.ValueOf(num)
	//v.Int()
	// 创建 UserService 实例
	userSvc := new(Server.UserService)

	// 创建 RPC 服务器
	rpcServer := Server.NewServer()
	// 注册 UserService
	err := rpcServer.RegisterService(userSvc)

	if err != nil {
		log.Fatalf("failed to register service: %v", err)
	}

	fmt.Println("Starting RPC server on :8080...")
	// 启动服务器
	if err := rpcServer.Start(":8080"); err != nil {
		log.Fatalf("server start error: %v", err)
	}
}
