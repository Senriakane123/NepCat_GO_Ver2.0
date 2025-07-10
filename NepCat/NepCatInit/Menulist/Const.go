package Menulist

var serverMenu = []string{
	"请选择你的服务",
	"请求格式为发送@机器人+对应服务字符即可使用符号或是跳转二级或者三级菜单",
	"例如发送'@机器人群管理'，则会返回对应二级菜单，亦或是启动某项服务",
	"1 群管理",
	"2 随机涩图 ",
	"3 更换机器人头像（需要向机器人所有者获取管理权限，目前是测试开发阶段请求格式为'@Bot修改头像-图片url'） ",
	"4 宠物系统（测试开发阶段）",
	"5 随机音乐推荐（'@Bot随机音乐推荐'，还在完善开发中）",
}

var GroupManageServerMenu = []string{
	"tips：此服务需要拥有管理员权限",
	"禁言类管理以分钟为计算，如果填入60则禁言1小时",
	"1 禁言 禁言格式为 '@需要禁言群友60",
	"例如发送'@sachiko60',则会禁言sachiko60分钟",
	"2 全体禁言 格式为 ‘@Bot全体禁言’",
	"3 解除全体禁言 格式为 ‘@Bot解除全体禁言’",
	"4 踢人 格式为'@群友踢人'",
	"例如发送'@sachiko踢人',则会把sachiko踢出群聊",
}

func GetServerList() []string {
	return serverMenu
}

func GetGroupServerList() []string {
	return GroupManageServerMenu
}
