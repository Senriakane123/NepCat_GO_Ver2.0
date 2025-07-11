package Handle

import (
	"NepCat_GO/NepCatInit/Menulist"
	"NepCat_GO/Tool"
)

func MenuReplyMsgBuild(qqnum string) string {
	return Tool.BuildAtQQString(qqnum) + "\n" + Tool.BuildReplyMessage(Menulist.GetServerList())
}

func GroupManageReplyMsgBuild(qqnum string) string {
	return Tool.BuildAtQQString(qqnum) + "\n" + Tool.BuildReplyMessage(Menulist.GetGroupServerList())
}

func PicReplyMsgBuild(qqnum string) string {
	return Tool.BuildAtQQString(qqnum) + "\n" + Tool.BuildReplyMessage(Menulist.GetPicServerList())
}
