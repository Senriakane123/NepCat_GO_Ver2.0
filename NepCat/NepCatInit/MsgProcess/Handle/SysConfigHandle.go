package Handle

import (
	ConfigManage "NepCat_GO/ConfigModule"
	"strings"
)

func ChangeReplayMode(rawmsg string) bool {
	//if !strings.Contains(rawmsg, "切换回复模式") {
	//	return false
	//}

	allowedModes := map[string]bool{
		"全回复":   true,
		"部分回复":  true,
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
