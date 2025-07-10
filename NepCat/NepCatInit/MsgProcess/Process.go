package MsgProcess

import (
	ConfigManage "NepCat_GO/ConfigModule"
	"NepCat_GO/NepCatInit/MSGModel"
	"NepCat_GO/NepCatInit/Menulist"
	"NepCat_GO/Tool"
	"fmt"
	"strconv"
	"strings"
)

func MessageRrocess(message MSGModel.ResMessage) {
	//获取原消息包含的所有QQ号
	for {
		_, QQNumberList := Tool.ListQQNumber(message.RawMessage)
		BContainBotQQ := true
		//检查是否包含原机器人QQ
		for _, n := range QQNumberList {
			if n == strconv.Itoa(int(message.SelfID)) {
				BContainBotQQ = true
			} else {
				BContainBotQQ = false
			}
		}
		fmt.Println(BContainBotQQ)
		fmt.Println("------------------此回复不包含QQ机器人账号,包含的账号列表为：", QQNumberList)

		switch ConfigManage.GetWebConfig().Mode.ReplyMode {
		case "全回复":
			//判断是否包含@机器请求的两种情况
			if BContainBotQQ {
				//切换内存回复模式
				if ChangeReplayMode(message.RawMessage) {
					ReplyNormalGroupMsg(message, "切换回复模式成功")
					break
				}
				//判断是否包含菜单请求
				if strings.Contains(message.RawMessage, "菜单") {
					ReplyNormalGroupMsg(message, MenuReplyMsgBuild(strconv.Itoa(int(message.SelfID))))
				}
			} else {

			}
			//fmt.Println("此项为发送过来的源消息", message.RawMessage)
		case "部分回复":
		case "管理员回复":
		case "开发者回复":
		}

		break
	}

	return

}

func ChangeReplayMode(rawmsg string) bool {
	if !strings.Contains(rawmsg, "切换回复模式") {
		return false
	}

	allowedModes := map[string]bool{
		"全回复":     true,
		"部分回复":   true,
		"管理员回复": true,
		"开发者回复": true,
	}

	for msg, _ := range allowedModes {
		if strings.Contains(rawmsg, msg) {
			ConfigManage.GetWebConfig().Mode.ReplyMode = msg
			return true
		}
	}
	return false
}

func MenuReplyMsgBuild(qqnum string) string {
	return Tool.BuildAtQQString(qqnum) + Tool.BuildReplyMessage(Menulist.GetServerList())
}
