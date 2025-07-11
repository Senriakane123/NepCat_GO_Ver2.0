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

var MenuCommand = map[string]func(MSGModel.ResMessage){}

func MenuInit() {
	MenuCommand = map[string]func(MSGModel.ResMessage){
		"菜单": func(msg MSGModel.ResMessage) {
			reply := Handle.MenuReplyMsgBuild(strconv.Itoa(int(msg.SelfID)))
			Handle.ReplyNormalGroupMsg(msg, reply)
		},
		"服务器状态": func(msg MSGModel.ResMessage) {
			reply := Handle.ServerStatusBuild(strconv.Itoa(int(msg.SelfID)))
			Handle.ReplyNormalGroupMsg(msg, reply)
		},
		"涩图管理": func(msg MSGModel.ResMessage) {
			reply := Handle.PicReplyMsgBuild(strconv.Itoa(int(msg.SelfID)))
			Handle.ReplyNormalGroupMsg(msg, reply)
		},
		"群管理": func(msg MSGModel.ResMessage) {
			reply := Handle.GroupManageReplyMsgBuild(strconv.Itoa(int(msg.SelfID)))
			Handle.ReplyNormalGroupMsg(msg, reply)
		},
		"切换回复模式": func(msg MSGModel.ResMessage) {
			Handle.ChangeReplayMode(msg.RawMessage)
			Handle.ReplyNormalGroupMsg(msg, "切换回复模式成功")
		},
	}
}

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

				for _, n := range GroupKeyWord {
					if strings.Contains(message.RawMessage, n) {
						Handle.GroupManage(message)
					}
				}

				for _, n := range PicKeyWord {
					if strings.Contains(message.RawMessage, n) {
						//Handle.GroupManage(message)
						Handle.RandomPicManage(message)
					}
				}

			}

			//判断是否包含@机器请求的两种情况
			if BContainBotQQ {
				for cmd, handler := range MenuCommand {
					if strings.Contains(message.RawMessage, cmd) {
						handler(message)
						break
					}
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
