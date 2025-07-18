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
	if offset <= 0 {
		return fmt.Errorf("解析响应协议头失败")
	}

	var respRpcHeader VRTSProxyRPCHeader
	ret := respRpcHeader.Parse(resp[offset:])
	if ret <= 0 {
		return fmt.Errorf("解析响应RPC头失败")
	}

	if respHeader.msgSn != header.msgSn {
		return fmt.Errorf("SN 不匹配，期待: %d, 实际: %d", header.msgSn, respHeader.msgSn)
	}

	fmt.Println("注册成功，收到 ServerID:", respRpcHeader.Service)
	return nil
}
