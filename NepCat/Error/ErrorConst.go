package Error

import "sync"

type ErrorConst struct {
	Errcode int
	ErrMsg  string
}

var (
	errorCodeMap map[int]string
	once         sync.Once
)

const (
	NEPCAT_NORLMAL_SUCCESS = 0
)

// 初始化函数（只执行一次）
func initErrorMap() {
	errorCodeMap = map[int]string{
		NEPCAT_NORLMAL_SUCCESS: "操作成功",
		// 可继续扩展...
	}
}

// 外部调用接口：根据错误码返回结构体（自动初始化）
func GetError(code int) ErrorConst {
	// 确保只初始化一次
	once.Do(initErrorMap)

	msg, ok := errorCodeMap[code]
	if !ok {
		msg = "未定义错误码"
	}
	return ErrorConst{
		Errcode: code,
		ErrMsg:  msg,
	}
}
