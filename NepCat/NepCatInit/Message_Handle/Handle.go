package Message_Handle

import (
	"NepCat_GO/NepCatInit/MSGModel"
	"NepCat_GO/NepCatInit/MsgProcess"
	"NepCat_GO/NepCatInit/Nepcat_ws_init"
	"encoding/json"
	"fmt"
)

func MessageHandler() {

	for msg := range *Nepcat_ws_init.NepcatWS.GetChannel() {
		var resmsg MSGModel.ResMessage
		err := json.Unmarshal([]byte(msg), &resmsg)
		if err != nil {
			fmt.Println("解析 JSON 失败:", err)
			return
		}

		switch resmsg.PostType {
		case "meta_event":
			handleMetaEvent(resmsg)
		case "message":
			handleChatMessage(resmsg)
		default:
			fmt.Println("未知消息类型:", resmsg.PostType)
		}
		//return
	}

}

func handleMetaEvent(msg MSGModel.ResMessage) {
	if msg.MetaEventType == "heartbeat" {
		fmt.Println("收到心跳包")
	} else if msg.MetaEventType == "lifecycle" {
		fmt.Println("机器人上线")
	}
}

func handleChatMessage(msg MSGModel.ResMessage) {
	fmt.Printf("收到来自用户 %d 的消息: %s\n", msg.UserID, msg.RawMessage)
	switch msg.MessageType {
	case "group":
		go MsgProcess.MessageRrocess(msg)
		fmt.Println("暂定调用http回复")
	case "private":
		//DS链接
		//go func() {
		//	var DSHandle DeepSeekReqHandle.DeepSeekManageHandle
		//	DSHandle.HandlerInit()
		//	if DSHandle.HandleGroupManageMessage(msg) {
		//		return
		//	}
		//}()

	}
	//if msg.MessageType == "group" {
	//	fmt.Printf("群聊消息，群ID: %d\n", msg.GroupID)
	//	// 可进一步处理 at 消息、图片等
	//}
}
