package Message_Handle

import (
	"NepCat_GO/NepCatInit/MSGModel"
	"NepCat_GO/NepCatInit/MsgProcess"
	"NepCat_GO/NepCatInit/Nepcat_ws_init"
	"encoding/json"
	"fmt"
	"github.com/jander/golog/logger"
)

func MessageHandler() {

	ws := Nepcat_ws_init.NepcatWS
	if ws == nil {
		fmt.Println("❌ WebSocket 未初始化，无法处理消息")
		return
	}

	ch := ws.GetChannel()
	if ch == nil {
		fmt.Println("❌ WebSocket Channel 未初始化")
		return
	}

	for msg := range *ch {
		var resmsg MSGModel.ResMessage
		err := json.Unmarshal([]byte(msg), &resmsg)
		if err != nil {
			fmt.Println("解析 JSON 失败:", err)
			continue
		}

		switch resmsg.PostType {
		case "meta_event":
			handleMetaEvent(resmsg)
		case "message":
			handleChatMessage(resmsg)
		default:
			fmt.Println("未知消息类型:", resmsg.PostType)
		}
	}
}

func handleMetaEvent(msg MSGModel.ResMessage) {
	if msg.MetaEventType == "heartbeat" {
		fmt.Println("收到心跳包")
	} else if msg.MetaEventType == "lifecycle" {
		logger.Info("机器人上线")
	}
}

func handleChatMessage(msg MSGModel.ResMessage) {
	fmt.Printf("收到来自用户 %d 的消息: %s\n", msg.UserID, msg.RawMessage)
	switch msg.MessageType {
	case "group":
		logger.Info("开启携程进行消息处理")
		logger.Info("处理消息的GroupID为：", msg.GroupID)
		go MsgProcess.MessageRrocess(msg)
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
