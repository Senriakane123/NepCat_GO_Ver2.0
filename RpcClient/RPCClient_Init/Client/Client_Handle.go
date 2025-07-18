package Client

import (
	"RpcClient/RPCClient_Init/Const"
	"RpcClient/Tool"
	"bytes"
	"fmt"
)

func (obj *CLient) RpcServer_Register(serverType string) error {
	var err error

	// 构建 RPC 注册消息
	var header VRTSProxyProtocolHeader
	var regMsg VRTSProxyRPCHeader

	// 填写协议头
	header.msgType = Const.VRTS_COMANND_TYPE_REGIST
	header.msgSn = uint32(Tool.GenerateNumSN()) // 随机 SN
	regMsg.ServerType = serverType
	regMsg.Method = 0   // 如有特定注册方法码，可定义
	regMsg.BodySize = 0 // 注册无消息体
	regMsg.Host = ""    // 可选，也可留空
	regMsg.Caller = 0
	regMsg.Service = 0

	// 设置包体长度
	header.size = regMsg.Size()

	// 构建完整包
	var buf bytes.Buffer
	header.Package(&buf)
	regMsg.Package(&buf)

	// 发送到服务端
	_, err = obj.Conn.Write(buf.Bytes())
	if err != nil {
		fmt.Println("发送注册数据失败:", err)
		return err
	}

	fmt.Println("注册包发送成功，ServerType:", serverType)

	// ✅ 等待响应（阻塞）
	resp := make([]byte, 1024)
	_, err = obj.Conn.Read(resp)
	if err != nil {
		fmt.Println("接收注册响应失败:", err)
		return err
	}

	// 解析响应头和RPC头
	var respHeader VRTSProxyProtocolHeader
	offset := respHeader.Parse(resp)
	if offset < 0 {
		return fmt.Errorf("解析响应协议头失败")
	}

	var respRpcHeader VRTSProxyRPCHeader
	ret := respRpcHeader.Parse(resp[respHeader.Size():])
	if ret < 0 {
		return fmt.Errorf("解析响应RPC头失败")
	}

	if respHeader.msgSn != header.msgSn {
		return fmt.Errorf("SN 不匹配，期待: %d, 实际: %d", header.msgSn, respHeader.msgSn)
	}

	fmt.Println(respHeader, respRpcHeader)
	fmt.Println("注册成功，收到 ServerID:", respRpcHeader.Service)
	obj.ServiceID = int64(respRpcHeader.Service)
	return nil
}

func (obj *CLient) ListenAndHandleServerMessages() {
	buf := make([]byte, 4096)

	for {
		n, err := obj.Conn.Read(buf)
		if err != nil {
			fmt.Printf("接收消息错误: %v\n", err)
			break
		}

		if n == 0 {
			fmt.Println("服务端关闭连接")
			break
		}

		// 解析协议头
		var header VRTSProxyProtocolHeader
		offset := header.Parse(buf[:n])
		if offset <= 0 {
			fmt.Println("解析协议头失败")
			continue
		}

		// 解析 RPC 头
		var rpcHeader VRTSProxyRPCHeader
		ret := rpcHeader.Parse(buf[offset:n])
		if ret <= 0 {
			fmt.Println("解析RPC头失败")
			continue
		}

		// 处理消息体（如果有）
		bodyStart := offset + int(ret)
		var msgBody string
		if bodyStart < n {
			msgBody = string(buf[bodyStart:n])
		}

		// 调用你自定义的处理逻辑
		fmt.Printf("收到服务端消息: Method=%d ServerType=%s Msg=%s\n", rpcHeader.Method, rpcHeader.ServerType, msgBody)

		// 可以根据 Method 做不同处理，如回调函数、命令派发等
	}
}
