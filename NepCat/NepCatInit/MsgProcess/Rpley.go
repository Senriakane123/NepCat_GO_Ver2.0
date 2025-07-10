package MsgProcess

import (
	"NepCat_GO/NepCatInit"
	"NepCat_GO/NepCatInit/MSGModel"
)

var RespApiMap = make(map[string]func(MSGModel.ResMessage, string))

func RespApiInit() {
	RespApiMap[NepCatInit.SEND_GROUP_MSG] = ReplyNormalGroupMsg
}

func ReplyNormalGroupMsg(rawmsg MSGModel.ResMessage, RespMsg string) {
	// 构造返回数据
	message := map[string]interface{}{
		"group_id": rawmsg.GroupID, // 群号
		"message":  RespMsg,
	}

	if handler, exists := NepCatInit.ReqApiMap[NepCatInit.SEND_GROUP_MSG]; exists {
		handler(NepCatInit.SEND_GROUP_MSG, message)
	}

}
