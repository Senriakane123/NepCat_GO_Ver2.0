package Handle

import (
	"NepCat_GO/NepCatInit/MSGModel"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func RandomPicManage(rawmsg MSGModel.ResMessage) (isreply bool) {
	fmt.Println("原始消息：", rawmsg.RawMessage)

	var qq string
	var num int = 1 // 默认图数为 1
	var tags []string
	isreply = false

	// 1️⃣ 提取 QQ 号
	qqRe := regexp.MustCompile(`qq=(\d+)`)
	qqMatches := qqRe.FindStringSubmatch(rawmsg.RawMessage)
	if len(qqMatches) >= 2 {
		qq = qqMatches[1]
		fmt.Println("✅ 提取到 QQ:", qq)
	} else {
		fmt.Println("⚠️ 未找到 QQ")
		return
	}

	// 2️⃣ 提取 CQ码后文本部分
	parts := strings.SplitN(rawmsg.RawMessage, "]", 2)
	if len(parts) < 2 {
		fmt.Println("⚠️ 消息格式不正确")
		return
	}
	text := strings.TrimSpace(parts[1])

	// 3️⃣ 判断是否是“随机涩图”开头
	if strings.HasPrefix(text, "随机涩图") {
		text = strings.TrimPrefix(text, "随机涩图")
		if strings.HasPrefix(text, "-") {
			text = strings.TrimPrefix(text, "-")
			segments := strings.SplitN(text, "-", 2)

			// 数量
			numStr := segments[0]
			if n, err := strconv.Atoi(numStr); err == nil {
				num = n
			}

			// 标签
			if len(segments) > 1 {
				tags = strings.Split(segments[1], "，") // 使用中文逗号
			}
		}
	}

	fmt.Println("🎯 QQ号:", qq)
	fmt.Println("🎯 数量:", num)
	fmt.Println("🎯 标签:", tags)

	return
}
