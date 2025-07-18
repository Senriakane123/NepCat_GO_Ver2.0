package Server

import (
	"RPCServer/Const"
	"RPCServer/Tool"
	"bytes"
	"fmt"
)

func (obj *Server) HandleConnection() {
	defer func() {
		// 清理连接
		addr := obj.Conn.RemoteAddr()
		fmt.Println("RPC Server: 客户端断开连接:", addr)
		Init_Server.Services.Delete(addr)
		obj.Conn.Close()
	}()

	for {
		msg, err, offset := obj.ReceivceMsg()
		if err != nil || offset == 0 {
			fmt.Println("接收消息错误:", err)
			return
		}

		var header VRTSProxyProtocolHeader
		header.Parse(msg)

		// 检查 SN 是否重复（可选，简单判断）
		if Tool.Contains(obj.MsgSnCodeList, int(header.msgSn)) {
			fmt.Println("重复 SN 消息, 忽略:", header.msgSn)
			continue
		}
		obj.MsgSnCodeList = append(obj.MsgSnCodeList, int(header.msgSn))

		// 根据消息类型分发处理
		switch header.msgType {
		case Const.VRTS_COMMAND_TYPE_HEARTBEAT:
			fmt.Println("收到心跳包")
			obj.ReplyHeartbeat(header)
		case Const.VRTS_COMMNAND_TYPE_CALL:
			obj.HandleCall(msg, offset, header)
		case Const.VRTS_COMMAND_RESP_MASK:
			fmt.Println("收到客户端响应，跳过")
		case Const.VRTS_COMANND_TYPE_REGIST:
			obj.ResgisterServer(msg, offset, header)
		default:
			fmt.Println("未知消息类型:", header.msgType)
		}
	}
}

func (obj *Server) ReceivceMsg() ([]byte, error, int) {
	var header VRTSProxyProtocolHeader
	offset := 0
	rcv_size := header.Size() // 通常为 16
	var err error
	var n int
	rcv_header := false

	msg := make([]byte, 4096) // ⭐ 初始化避免 panic，可按需扩大

	for {
		n, err = obj.Conn.Read(msg[offset:rcv_size])
		if err != nil {
			fmt.Println("接收失败:", err)
			break
		}
		if n <= 0 {
			break
		}
		offset += n

		if offset < header.Size() {
			continue
		}

		if !rcv_header {
			header.Parse(msg)
			rcv_header = true
		}

		rcv_size = int(header.size) + header.Size()
		if offset >= rcv_size {
			break
		}
	}

	if n > 0 {
		return msg[:rcv_size], nil, rcv_size // 🔒 截取实际消息长度
	}

	return nil, err, 0
}

func (obj *Server) ReplyHeartbeat(header VRTSProxyProtocolHeader) {
	var buf bytes.Buffer
	header.msgType = Const.VRTS_COMMAND_TYPE_HEARTBEAT // 重复使用原 SN 和版本
	header.size = 0
	header.Package(&buf)
	obj.Conn.Write(buf.Bytes())
	fmt.Println("已回复心跳")
}

func (obj *Server) HandleCall(rcvBuf []byte, Len int, header VRTSProxyProtocolHeader) {
	var rpcHeader VRTSProxyRPCHeader
	retVal := rpcHeader.Parse(rcvBuf[header.Size():])
	if retVal <= 0 {
		fmt.Println("Receive invalid rpc call")
		return
	}

	msgBody := rcvBuf[header.Size()+int(retVal) : Len]

	// 遍历查找目标服务
	var targetServer *Server
	Init_Server.Services.Range(func(key, value interface{}) bool {
		srv := value.(*Server)
		if srv.ServerType == rpcHeader.ServerType {
			targetServer = srv
			return false // 找到就停止遍历
		}
		return true
	})

	if targetServer == nil {
		fmt.Println("未找到匹配的服务类型:", rpcHeader.ServerType)
		return
	}

	// 构建转发请求包
	var forwardBuf bytes.Buffer
	forwardHeader := VRTSProxyProtocolHeader{
		version: 1,
		msgType: Const.VRTS_COMMNAND_TYPE_CALL,
		msgSn:   header.msgSn, // 可复用原始 SN
		size:    int32(len(msgBody)) + rpcHeader.Size(),
	}
	forwardHeader.Package(&forwardBuf)
	rpcHeader.Package(&forwardBuf)
	forwardBuf.Write(msgBody)

	// 发送给目标服务
	_, err := targetServer.Conn.Write(forwardBuf.Bytes())
	if err != nil {
		fmt.Println("转发请求失败:", err)
	}

	// ✅ 等待目标客户端返回响应
	response := make([]byte, 4096) // 你可根据预期消息大小自行调整
	n, err := targetServer.Conn.Read(response)
	if err != nil {
		fmt.Println("读取目标客户端响应失败:", err)
		return
	}

	fmt.Println("收到目标客户端响应:", string(response[:n]))
	// 👇你可自行根据收到的响应内容转发回原始请求方
}

func (obj *Server) ResgisterServer(rcvBuf []byte, len int, header VRTSProxyProtocolHeader) {
	var rpcHeader VRTSProxyRPCHeader
	retVal := rpcHeader.Parse(rcvBuf[header.Size():])
	if retVal > 0 {
		msgBody := string(rcvBuf[header.Size()+int(retVal) : len])
		//obj.handle.DoCommand(rpcHeader.Method, msgBody)
		fmt.Println(msgBody)
	} else {
		fmt.Println("Length:", len)
		fmt.Println("PackageHeaderSize:", header.Size())
		fmt.Println("Receive invalid rpc call")
	}
	obj.ServerType = rpcHeader.ServerType

	var respHeader VRTSProxyProtocolHeader
	var respRpcHeader VRTSProxyRPCHeader

	respHeader.version = 1
	respHeader.msgType = Const.VRTS_COMMAND_RESP_MASK // 假设常量定义为 ACK 类型
	respHeader.msgSn = header.msgSn                   // 保持 SN 一致

	respRpcHeader.Method = 0
	respRpcHeader.ServerType = obj.ServerType
	respRpcHeader.Host = ""
	respRpcHeader.Caller = 0
	respRpcHeader.Service = int32(obj.ServiceID) // 将 ServerID 填进响应里
	respRpcHeader.BodySize = 0                   // 无消息体

	respHeader.size = respRpcHeader.Size()

	var respBuf bytes.Buffer
	fmt.Println(respHeader, respRpcHeader)
	respHeader.Package(&respBuf)
	respRpcHeader.Package(&respBuf)

	_, err := obj.Conn.Write(respBuf.Bytes())
	if err != nil {
		fmt.Println("注册响应发送失败:", err)
	}
}
