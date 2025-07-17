package Server

import (
	"fmt"
	RPCCommon "github.com/gopublic/RPCStruct"
	"net"
	"reflect"
	"sync"
)

func init() {
	fmt.Println(RPCCommon.CompressTypeNone)
}

// ServiceMethod 结构体存储服务方法信息
type ServiceMethod struct {
	Method     reflect.Value // 方法的反射值
	MethodType reflect.Type  // 方法的类型
}

// Service 结构体存储服务信息
type Service struct {
	Name    string                    // 服务名
	Rcvr    reflect.Value             // 服务的实例
	Methods map[string]*ServiceMethod // 服务方法集合
}

// Server RPC服务器
type Server struct {
	listener net.Listener
	services sync.Map // 存储注册的服务: map[string]*Service
	seq      uint64   // 请求序列号
}
