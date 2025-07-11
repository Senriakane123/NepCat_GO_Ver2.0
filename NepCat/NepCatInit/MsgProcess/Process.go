package MsgProcess

import (
	ConfigManage "NepCat_GO/ConfigModule"
	"NepCat_GO/NepCatInit/MSGModel"
	"NepCat_GO/NepCatInit/MsgProcess/Handle"
	"NepCat_GO/Tool"
	"fmt"
	"strconv"
	"strings"
)

func MessageRrocess(message MSGModel.ResMessage) {
	//获取原消息包含的所有QQ号
	for {
		_, QQNumberList := Tool.ListQQNumber(message.RawMessage)
		BContainBotQQ := false
		//检查是否包含原机器人QQ
		for _, n := range QQNumberList {
			if n == strconv.Itoa(int(message.SelfID)) {
				BContainBotQQ = true
				break
			}
		}
		fmt.Println(BContainBotQQ)
		fmt.Println("------------------此回复不包含QQ机器人账号,包含的账号列表为：", QQNumberList)

		switch ConfigManage.GetWebConfig().Mode.ReplyMode {
		case "全回复":
			//处理群管理消息
			if len(QQNumberList) != 0 {
				isreply := Handle.GroupManage(message)
				if isreply == true {
					break
				}
			}

			//判断是否包含@机器请求的两种情况
			if BContainBotQQ {
				//切换内存回复模式
				if strings.Contains(message.RawMessage, "切换回复模式") {
					Handle.ChangeReplayMode(message.RawMessage)
					Handle.ReplyNormalGroupMsg(message, "切换回复模式成功")
					break
				}
				//判断是否包含菜单请求
				if strings.Contains(message.RawMessage, "菜单") {
					Handle.ReplyNormalGroupMsg(message, Handle.MenuReplyMsgBuild(strconv.Itoa(int(message.SelfID))))
					break
				}
				//获取服务器状态
				if strings.Contains(message.RawMessage, "服务器状态") {
					Handle.ReplyNormalGroupMsg(message, Handle.ServerStatusBuild(strconv.Itoa(int(message.SelfID))))
					break
				}

				if strings.Contains(message.RawMessage, "涩图管理") {
					Handle.ReplyNormalGroupMsg(message, Handle.PicReplyMsgBuild(strconv.Itoa(int(message.SelfID))))
					break
				}

				if strings.Contains(message.RawMessage, "群管理") {
					Handle.ReplyNormalGroupMsg(message, Handle.GroupManageReplyMsgBuild(strconv.Itoa(int(message.SelfID))))
					break
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
	//ReplyNormalGroupMsg(message, "不管你是谁如果看到这条信息请让制作组滚去敲代码")

	return

}
