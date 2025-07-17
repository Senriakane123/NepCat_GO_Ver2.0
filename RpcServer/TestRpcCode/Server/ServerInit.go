package Server

import (
	"bytes" // 导入 bytes 包
	"encoding/gob"
	"fmt"
	RPCCommon "github.com/gopublic/RPCStruct" // 替换为你的项目路径
	"io"
	"net"
	"reflect"
)

// NewServer 创建一个新的RPC服务器
func NewServer() *Server {
	return &Server{}
}

// RegisterService 注册服务
// rcvr 必须是一个可导出类型（首字母大写）的实例或指针
func (s *Server) RegisterService(rcvr interface{}) error {
	service := new(Service)
	service.Rcvr = reflect.ValueOf(rcvr)
	//service.Rcvr.MethodByName("")
	// 获取服务名，通常是结构体的类型名
	service.Name = reflect.TypeOf(rcvr).Elem().Name()
	if service.Name == "" {
		return fmt.Errorf("rpc server: no service name in type %T", rcvr)
	}

	service.Methods = suitableMethods(reflect.TypeOf(rcvr))

	if len(service.Methods) == 0 {
		return fmt.Errorf("rpc server: type %T has no exported methods of suitable type", rcvr)
	}

	s.services.Store(service.Name, service)
	fmt.Printf("RPC server: registered service %s with %d methods\n", service.Name, len(service.Methods))
	return nil
}

// suitableMethods 查找服务中可导出的RPC方法
// 方法签名约定： func (receiver *T) Method(args *ArgsType) (reply *ReplyType, err error)
func suitableMethods(typ reflect.Type) map[string]*ServiceMethod {
	methods := make(map[string]*ServiceMethod)
	for i := 0; i < typ.NumMethod(); i++ {
		method := typ.Method(i)
		mType := method.Type
		// 方法必须是可导出的 (首字母大写)
		if method.PkgPath != "" { // PkgPath 非空表示方法是未导出的
			continue
		}
		// 方法必须至少有两个参数 (receiver, args)
		// 方法必须返回两个值 (reply, error)
		if mType.NumIn() != 2 || mType.NumOut() != 2 {
			continue
		}
		// 第二个返回值必须是error类型
		if mType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		// 第一个参数 (args) 必须是可导出的结构体或指针
		argType := mType.In(1)
		if argType.Kind() != reflect.Ptr && argType.Kind() != reflect.Struct {
			continue
		}
		// 第一个返回值 (reply) 必须是指针
		replyType := mType.Out(0)
		if replyType.Kind() != reflect.Ptr {
			continue
		}
		methods[method.Name] = &ServiceMethod{Method: method.Func, MethodType: mType}
	}
	return methods
}

// Start 启动RPC服务器
func (s *Server) Start(addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("rpc server: failed to listen on %s: %v", addr, err)
	}
	s.listener = listener
	fmt.Printf("RPC server: listening on %s\n", addr)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("RPC server: accept error: %v\n", err)
			continue
		}
		go s.handleConnection(conn)
	}
}

// handleConnection 处理单个客户端连接
func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	// 创建一个用于从连接读取数据的解码器
	decoder := gob.NewDecoder(conn)
	// 创建一个用于向连接写入数据的编码器 (用于消息头)
	encoder := gob.NewEncoder(conn)

	for {
		// 1. 读取消息头
		var msgHeader RPCCommon.RPCMessageHeader
		err := decoder.Decode(&msgHeader)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("RPC server: failed to decode message header: %v\n", err)
			}
			return // 连接关闭或发生错误，终止处理
		}
		fmt.Printf("RPC server: received message header: %s\n", &msgHeader)

		// 检查魔数
		if msgHeader.MagicNumber != RPCCommon.MagicNumber {
			fmt.Printf("RPC server: invalid magic number: %x\n", msgHeader.MagicNumber)
			return
		}

		// 检查消息类型 (只处理请求)
		if msgHeader.MessageType != RPCCommon.MessageTypeRequest {
			fmt.Printf("RPC server: unexpected message type: %d\n", msgHeader.MessageType)
			return
		}

		// 2. 读取消息体 (请求头 + 请求参数)
		payloadBytes := make([]byte, msgHeader.BodyLength)
		_, err = io.ReadFull(conn, payloadBytes) // 从连接中精确读取指定长度的字节
		if err != nil {
			fmt.Printf("RPC server: failed to read payload bytes: %v\n", err)
			return
		}

		// 使用新的解码器处理消息体
		payloadDecoder := gob.NewDecoder(bytes.NewReader(payloadBytes))

		// 3. 读取请求头
		var reqHeader RPCCommon.RPCRequestHeader
		err = payloadDecoder.Decode(&reqHeader)
		if err != nil {
			fmt.Printf("RPC server: failed to decode request header from payload: %v\n", err)
			s.sendResponse(encoder, conn, &reqHeader, nil, fmt.Errorf("failed to decode request header: %v", err))
			continue
		}
		fmt.Printf("RPC server: received request header: %s\n", &reqHeader)

		// 4. 查找服务和方法
		serviceVal, ok := s.services.Load(reqHeader.ServiceID)
		if !ok {
			s.sendResponse(encoder, conn, &reqHeader, nil, fmt.Errorf("service %s not found", reqHeader.ServiceID))
			continue
		}
		service := serviceVal.(*Service)

		method, ok := service.Methods[reqHeader.MethodName]
		if !ok {
			s.sendResponse(encoder, conn, &reqHeader, nil, fmt.Errorf("method %s not found in service %s", reqHeader.MethodName, reqHeader.ServiceID))
			continue
		}

		// 5. 反序列化请求参数
		argsType := method.MethodType.In(1) // 获取方法的第二个参数类型 (请求参数)
		// reflect.New(argsType.Elem()) 创建一个指向 argsType 所指向的类型的指针
		args := reflect.New(argsType.Elem()).Interface() // 创建参数实例 (例如 *int, *User)
		err = payloadDecoder.Decode(args)                // 解码到参数实例
		if err != nil {
			s.sendResponse(encoder, conn, &reqHeader, nil, fmt.Errorf("failed to decode request arguments: %v", err))
			continue
		}

		// 6. 调用服务方法
		replyType := method.MethodType.Out(0) // 获取方法的第一个返回值类型 (响应结果)
		// reflect.New(replyType.Elem()) 创建一个指向 replyType 所指向的类型的指针
		reply := reflect.New(replyType.Elem()).Interface() // 创建响应实例 (例如 *User)

		// 调用方法，传递服务的接收者和参数
		returns := method.Method.Call([]reflect.Value{service.Rcvr, reflect.ValueOf(args)})
		var callErr error
		if !returns[1].IsNil() { // 检查第二个返回值 (error)
			callErr = returns[1].Interface().(error)
		}

		// 7. 发送响应
		s.sendResponse(encoder, conn, &reqHeader, reply, callErr)
	}
}

// sendResponse 发送RPC响应
// encoder: 用于发送消息头的gob编码器
// conn: 底层网络连接
// reqHeader: 原始请求头，用于响应中的匹配
// reply: 实际的响应结果
// err: 调用过程中发生的错误
func (s *Server) sendResponse(encoder *gob.Encoder, conn net.Conn, reqHeader *RPCCommon.RPCRequestHeader, reply interface{}, err error) {
	respMsgHeader := RPCCommon.RPCMessageHeader{
		MagicNumber:     RPCCommon.MagicNumber,
		MessageType:     RPCCommon.MessageTypeResponse,
		CompressionType: RPCCommon.CompressTypeNone,
		SerializeType:   RPCCommon.SerializeTypeGob,
	}

	// 构造 RPCResponse 结构体
	rpcResponse := RPCCommon.RPCResponse{}
	if err != nil {
		rpcResponse.Error = err.Error() // 设置错误信息
	} else {
		rpcResponse.Reply = reply // 设置响应结果
	}

	// 1. 编码请求头和 RPCResponse 到缓冲区，以计算其总长度作为消息体长度
	var payloadBuf bytes.Buffer
	payloadEncoder := gob.NewEncoder(&payloadBuf)

	// 编码原始请求头 (用于客户端匹配响应)
	encodeErr := payloadEncoder.Encode(reqHeader)
	if encodeErr != nil {
		fmt.Printf("RPC server: failed to encode response request header to payload buffer: %v\n", encodeErr)
		return
	}

	// 编码 RPCResponse 结构体 (包含实际结果或错误)
	encodeErr = payloadEncoder.Encode(rpcResponse)
	if encodeErr != nil {
		fmt.Printf("RPC server: failed to encode RPCResponse to payload buffer: %v\n", encodeErr)
		return
	}

	respMsgHeader.BodyLength = uint32(payloadBuf.Len()) // 获取有效载荷的总长度

	// 2. 发送消息头
	encodeErr = encoder.Encode(&respMsgHeader) // 使用主编码器发送消息头
	if encodeErr != nil {
		fmt.Printf("RPC server: failed to encode response message header to connection: %v\n", encodeErr)
		return
	}

	// 3. 发送有效载荷 (请求头 + RPCResponse)
	_, encodeErr = conn.Write(payloadBuf.Bytes()) // 将缓冲区的内容直接写入连接
	if encodeErr != nil {
		fmt.Printf("RPC server: failed to write response payload to connection: %v\n", encodeErr)
		return
	}
}

// UserService 示例服务实现
type UserService struct{}

// GetUserById 根据ID获取用户信息
// 参数 id: 用户ID的指针
// 参数 user: 用于接收用户信息的指针
// 返回值 error: 如果发生错误则返回错误信息
func (u *UserService) GetUserById(id *int, user *RPCCommon.User) error {
	fmt.Printf("Server received GetUserById request with ID: %d\n", *id)
	// 模拟数据库查询
	if *id == 1 {
		user.ID = 1
		user.Name = "Alice"
		user.Email = "alice@example.com"
	} else if *id == 2 {
		user.ID = 2
		user.Name = "Bob"
		user.Email = "bob@example.com"
	} else {
		return fmt.Errorf("user with ID %d not found", *id)
	}
	return nil
}
