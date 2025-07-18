package Client

import (
	"fmt"
	"net"
)

func RpcClientMapInit() *Cient_Init {
	if &Init_Client == nil {
		Init_Client = Cient_Init{}
	} else {
		return &Init_Client
	}
	return &Init_Client
}

func (obj *Cient_Init) RpcClient_Init(addr string, servertypr string) error {
	Conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("rpc client: failed to connect on %s: %v", addr, err)
	}
	//obj.Listener = listener
	fmt.Printf("RPC server: connect on %s\n", addr)

	var NewClient = &CLient{
		//ServiceID:  int64(Init_Client_Num),
		Conn:       Conn,
		ServerType: servertypr,
		RemoteAddr: Conn.RemoteAddr().String(),
	}

	obj.Clients.Store(Conn.RemoteAddr(), NewClient)

	err = NewClient.RpcServer_Register(servertypr)
	if err != nil {
		fmt.Println("注册rpc服务失败：", err.Error())
	}
	// 注册成功后开启 goroutine 持续监听服务端消息
	go NewClient.ListenAndHandleServerMessages()

	return err

}
