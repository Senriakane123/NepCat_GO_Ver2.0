package Handle

import (
	"NepCat_GO/NepCatInit"
	"NepCat_GO/NepCatInit/MSGModel"
)

//var RespApiMap = make(map[string]func(MSGModel.ResMessage, string))
//
//func RespApiInit() {
//	RespApiMap[NepCatInit.SEND_GROUP_MSG] = ReplyNormalGroupMsg
//}

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

func ReplyBanMsg(groupid, usrid, time int64) {

	message := map[string]interface{}{
		"group_id": groupid, // 替换为你的QQ群号
		"user_id":  usrid,   // 发送的消息内容

		"duration": time, // 发送的消息内容
	}

	if handler, exists := NepCatInit.ReqApiMap[NepCatInit.SET_GROUP_BAN]; exists {
		handler(NepCatInit.SET_GROUP_BAN, message)
	}
}

func ReplyKickMsg(GroupId int64, UserId int64, kickboolen bool) {
	message := map[string]interface{}{
		"group_id": GroupId, // 替换为你的QQ群号
		"user_id":  UserId,  // 发送的消息内容

		"reject_add_request": kickboolen, // 发送的消息内容
	}

	if handler, exists := NepCatInit.ReqApiMap[NepCatInit.SET_GROUP_KICK]; exists {
		handler(NepCatInit.SET_GROUP_KICK, message)
	}
}

func ReplyGroupBanMsg(GroupId int64, BanBoolen bool) {
	message := map[string]interface{}{
		"group_id": GroupId,   // 替换为你的QQ群号
		"enable":   BanBoolen, // 发送的消息内容
	}

	if handler, exists := NepCatInit.ReqApiMap[NepCatInit.SET_GROUP_WHOLE_BAN]; exists {
		handler(NepCatInit.SET_GROUP_WHOLE_BAN, message)
	}
}
