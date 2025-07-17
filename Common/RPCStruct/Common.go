// common/protocol.go
package common

import (
	"encoding/gob"
	"fmt"
)

// MagicNumber 定义RPC消息的魔数，用于标识消息
const MagicNumber = [4]byte{'G', 'O', 'R', 'P'}

// MessageType 定义消息类型，例如请求或响应
type MessageType int

const (
	MessageTypeRequest  MessageType = iota // 请求消息
	MessageTypeResponse                    // 响应消息
)

// CompressionType 定义压缩类型
type CompressionType int

const (
	CompressTypeNone CompressionType = iota // 不压缩
	// CompressTypeGzip                   // Gzip压缩 (示例，可扩展)
)

// SerializeType 定义序列化类型
type SerializeType int

const (
	SerializeTypeGob  SerializeType = iota // Go内置的Gob编码
	SerializeTypeJSON                      // JSON编码 (示例，可扩展)
	// SerializeTypeProto                   // Protobuf编码 (示例，可扩展)
)

// RPCMessageHeader RPC消息的头部
// 包含魔数、消息类型、压缩类型、序列化类型、消息体长度
type RPCMessageHeader struct {
	MagicNumber     [4]byte         // 魔数，用于标识RPC消息
	MessageType     MessageType     // 消息类型
	CompressionType CompressionType // 压缩类型
	SerializeType   SerializeType   // 序列化类型
	BodyLength      uint32          // 消息体长度 (请求头 + 请求参数/响应结果 的总长度)
}

// RPCRequestHeader RPC请求的头部
// 包含服务名、方法名、请求ID
type RPCRequestHeader struct {
	ServiceID  string // 服务ID，例如 "UserService"
	MethodName string // 方法名，例如 "GetUserById"
	RequestID  uint64 // 请求ID，用于匹配请求和响应
}

// RPCResponse 结构体用于承载RPC响应
// 包含实际的响应数据和可能的错误信息
type RPCResponse struct {
	Reply interface{} // 实际的响应数据
	Error string      // 错误信息，如果调用成功则为空
}

// String 方法用于方便打印RPCMessageHeader
func (h *RPCMessageHeader) String() string {
	return fmt.Sprintf("Magic: %x, Type: %d, Compress: %d, Serialize: %d, BodyLen: %d",
		h.MagicNumber, h.MessageType, h.CompressionType, h.SerializeType, h.BodyLength)
}

// String 方法用于方便打印RPCRequestHeader
func (h *RPCRequestHeader) String() string {
	return fmt.Sprintf("Service: %s, Method: %s, ReqID: %d",
		h.ServiceID, h.MethodName, h.RequestID)
}

// User 示例用户结构体
type User struct {
	ID    int
	Name  string
	Email string
}

func init() {
	// 注册自定义类型，以便Gob能够正确编码和解码 interface{} 中的具体类型
	// 任何通过 interface{} 传递的自定义类型都需要在这里注册
	gob.Register(&User{})
	// gob.Register(User{}) // 也可以注册值类型，取决于实际使用情况
}
