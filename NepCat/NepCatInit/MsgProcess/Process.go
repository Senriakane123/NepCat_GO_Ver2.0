package MsgProcess

import (
	ConfigManage "NepCat_GO/ConfigModule"
	"NepCat_GO/NepCatInit/MSGModel"
	"NepCat_GO/Tool"
)

func MessageRrocess(message MSGModel.ResMessage) {
	//获取原消息包含的所有QQ号
	_, QQNumberList := Tool.ListQQNumber(message.RawMessage)
	switch ConfigManage.GetWebConfig().Mode.ReplyMode {
	case "全回复":
	case "部分回复":
	case "管理员回复":
	case "开发者回复":
	}
	for _, n := range QQNumberList {
		if n == string(message.SelfID) {

		}
	}
}
