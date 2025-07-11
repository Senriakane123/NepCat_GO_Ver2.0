package Tool

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// AtQQNumber 函数用于判断字符串中是否包含 @QQ 号，并返回包含的 @QQ 号列表
func ListQQNumber(str string) (bool, []string) {
	re := regexp.MustCompile(`\[CQ:at,qq=(\d+)\]`)
	matches := re.FindAllStringSubmatch(str, -1)
	var QQNumbers []string
	for _, match := range matches {
		if len(match) > 1 {
			QQNumbers = append(QQNumbers, match[1]) // 只存 QQ 号
		}
	}
	if len(QQNumbers) > 0 {
		return true, QQNumbers
	}
	return false, nil
}

func BuildAtQQString(QQnum string) string {
	return "[CQ:at,qq=" + QQnum + "]"
}

func BuildReplyMessage(Message []string) string {
	var builder strings.Builder

	for _, item := range Message {
		builder.WriteString(item)
		builder.WriteString("\n")
	}

	return builder.String()
}

func StringToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("无法转换字符串 '%s' 为整数: %v", s, err)
	}
	return i, nil
}

// int 转字符串
func IntToString(i int) string {
	return strconv.Itoa(i)
}
