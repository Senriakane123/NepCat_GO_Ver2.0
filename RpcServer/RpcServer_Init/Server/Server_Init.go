package Server

import (
	"fmt"
	"net"
)

func RpcServerMapInit() *Service_Init {
	if &Init_Server == nil {
		Init_Server = Service_Init{}
	} else {
		return &Init_Server
	}
	return &Init_Server
}

var Init_Server_Num = 1

func (obj *Service_Init) RpcServer_Init(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("rpc server: failed to listen on %s: %v", addr, err)
	}
	obj.Listener = listener
	fmt.Printf("RPC server: listening on %s\n", addr)

	for {
		conn, err := obj.Listener.Accept()
		if err != nil {
			fmt.Printf("RPC server: accept error: %v\n", err)
			continue
		} else {
			fmt.Printf("RPC server: accepted connection from %s\n", conn.RemoteAddr().String())
		}
		var NewServer = &Server{
			ServiceID:  int64(Init_Server_Num),
			Conn:       conn,
			RemoteAddr: conn.RemoteAddr().String(),
		}

		Init_Server_Num++
		obj.Services.Store(conn.RemoteAddr(), NewServer)
		go NewServer.HandleConnection()
	}

}
