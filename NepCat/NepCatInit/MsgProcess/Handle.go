package MsgProcess

import (
	ConfigManage "NepCat_GO/ConfigModule"
	"NepCat_GO/NepCatInit/Menulist"
	"NepCat_GO/SysStatusModule"
	"NepCat_GO/Tool"
	"fmt"
	"strings"
)

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
	return Tool.BuildAtQQString(qqnum) + "\n" + Tool.BuildReplyMessage(Menulist.GetServerList())
}

func GroupManageReplyMsgBuild(qqnum string) string {
	return Tool.BuildAtQQString(qqnum) + "\n" + Tool.BuildReplyMessage(Menulist.GetGroupServerList())
}

func PicReplyMsgBuild(qqnum string) string {
	return Tool.BuildAtQQString(qqnum) + "\n" + Tool.BuildReplyMessage(Menulist.GetPicServerList())
}

func ServerStatusBuild(qqnum string) string {
	status := SysStatusModule.GetSysInfo()
	var resp []string

	resp = append(resp, fmt.Sprintf("🕒 时间戳：%s", status.TimeStamp))
	resp = append(resp, fmt.Sprintf("🧠 CPU 使用率：%.2f%%", status.CPU))
	resp = append(resp, fmt.Sprintf("💾 内存使用率：%.2f%%", status.Memory))

	resp = append(resp, fmt.Sprintf("📶 网络接收速率：%.2f KB/s", status.Network.RecvBPS/1024))
	resp = append(resp, fmt.Sprintf("📤 网络发送速率：%.2f KB/s", status.Network.SendBPS/1024))

	resp = append(resp, fmt.Sprintf("📀 磁盘读取速率：%.2f KB/s", status.Disk.ReadBPS/1024))
	resp = append(resp, fmt.Sprintf("📝 磁盘写入速率：%.2f KB/s", status.Disk.WriteBPS/1024))

	return Tool.BuildAtQQString(qqnum) + "\n" + Tool.BuildReplyMessage(resp)
}
