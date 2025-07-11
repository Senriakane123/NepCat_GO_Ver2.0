package Handle

import (
	"NepCat_GO/NepCatInit/MSGModel"
	"NepCat_GO/Tool"
	"fmt"
	"regexp"
	"strings"
)

func GroupManage(rawmsg MSGModel.ResMessage) (isreply bool) {
	fmt.Println("原始消息：", rawmsg)

	var action string
	var value int
	var qq string
	var hasValue bool
	isreply = false

	// 1️⃣ 提取 QQ 号
	qqRe := regexp.MustCompile(`qq=(\d+)`)
	qqMatches := qqRe.FindStringSubmatch(rawmsg.RawMessage)
	if len(qqMatches) >= 2 {
		qq = qqMatches[1]
		fmt.Println("✅ 提取到 QQ:", qq)
	} else {
		fmt.Println("⚠️ 未找到 QQ")
	}

	// 2️⃣ 提取动作和参数
	// 假设动作在 CQ 码之后的部分（空格分隔）
	fields := strings.Fields(rawmsg.RawMessage)
	if len(fields) < 2 {
		fmt.Println("⚠️ 消息格式不正确")
		return
	}

	actionPart := fields[len(fields)-1]
	re := regexp.MustCompile(`^([^\d\s]+)(\d*)$`)
	matches := re.FindStringSubmatch(actionPart)
	if len(matches) >= 2 {
		action = matches[1]
		if len(matches[2]) > 0 {
			fmt.Sscanf(matches[2], "%d", &value)
			hasValue = true
		}
	}

	fmt.Printf("🎯 动作: %s\n", action)
	if hasValue {
		fmt.Printf("🎯 数值: %d\n", value)
	}
	qqInt, _ := Tool.StringToInt(qq)
	// 3️⃣ 执行命令逻辑
	switch action {
	case "禁言":
		if hasValue {
			fmt.Printf("👉 执行禁言 %s 用户 %d 秒\n", qq, value)
		} else {
			fmt.Println("⚠️ 禁言缺少时间参数")
		}
		ReplyBanMsg(rawmsg.GroupID, int64(qqInt), int64(value*60))
		isreply = true
	case "踢人":
		fmt.Printf("👉 执行踢出用户 %s\n", qq)
		ReplyKickMsg(rawmsg.GroupID, int64(qqInt), false)
		isreply = true
	case "全体禁言":
		fmt.Println("👉 执行全体禁言")
		ReplyGroupBanMsg(rawmsg.GroupID, true)
		isreply = true
	case "解除全体禁言":
		fmt.Println("👉 执行解除全体禁言")
		ReplyGroupBanMsg(rawmsg.GroupID, false)
		isreply = true
	default:
		fmt.Println("⚠️ 未知指令:", action)
	}
	return
}
