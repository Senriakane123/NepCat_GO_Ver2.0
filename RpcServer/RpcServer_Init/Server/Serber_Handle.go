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
	rcv_size := header.Size()
	var msg []byte
	var err error
	len := 0
	rcv_header := false

	for {

		len, err = obj.Conn.Read(msg[offset:rcv_size])
		if err != nil {
			fmt.Println(err.Error())
			break
		}

		if len < 0 {
			break
		}
		offset += len

		if offset < header.Size() {
			continue
		}

		if !rcv_header {
			header.Parse(msg)
			rcv_header = true
		}

		rcv_size = int(header.size) + header.Size()

		// 如果数据包接收完成
		if (header.Size() + int(header.size)) <= offset {
			break
		}

	}
	if len >= 0 {
		return msg, err, offset
	}

	return msg, err, 0
}

func (obj *Server) ReplyHeartbeat(header VRTSProxyProtocolHeader) {
	var buf bytes.Buffer
	header.msgType = Const.VRTS_COMMAND_TYPE_HEARTBEAT // 重复使用原 SN 和版本
	header.size = 0
	header.Package(&buf)
	obj.Conn.Write(buf.Bytes())
	fmt.Println("已回复心跳")
}

func (obj *Server) HandleCall(rcvBuf []byte, len int, header VRTSProxyProtocolHeader) {

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
}
