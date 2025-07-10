package MsgProcess

import (
	ConfigManage "NepCat_GO/ConfigModule"
	"NepCat_GO/NepCatInit/MSGModel"
	"NepCat_GO/Tool"
	"fmt"
)

func MessageRrocess(message MSGModel.ResMessage) {
	//获取原消息包含的所有QQ号
	_, QQNumberList := Tool.ListQQNumber(message.RawMessage)
	BContainBotQQ := true
	//检查是否包含原机器人QQ
	for _, n := range QQNumberList {
		if n == string(message.SelfID) {
			BContainBotQQ = true
		} else {
			BContainBotQQ = false
		}
	}
	fmt.Println(BContainBotQQ)
	fmt.Println("------------------此回复不包含QQ机器人账号,包含的账号列表为：", QQNumberList)
	switch ConfigManage.GetWebConfig().Mode.ReplyMode {
	case "全回复":
		ReplyNormalGroupMsg(message, "测试回复")
	case "部分回复":
	case "管理员回复":
	case "开发者回复":
	}

}
