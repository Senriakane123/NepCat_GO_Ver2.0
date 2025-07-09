package main

import (
	ConfigManage "NepCat_GO/ConfigModule"
	"NepCat_GO/NepCatInit"
	"NepCat_GO/NepCatInit/Nepcat_ws_init"
	db "github.com/gopublic/GormModule/DBControl/DatabaseControl"
	"github.com/jander/golog/logger"
)

func main() {
	rotatingHandler := logger.NewFileHandler("NepCat_RunningLog/Running.log")
	logger.SetHandlers(logger.Console, rotatingHandler)
	logger.Info("------------------------------------------------------------------------系统日志初始化------------------------------------------------------------------------")

	logger.Info("------------------------------------------------------------------------config初始化------------------------------------------------------------------------")
	err := ConfigManage.ConfigInit("ConfigModule/config.yaml") // 这里假设配置文件名为 config.yaml
	if err != nil {
		logger.Error("Failed to load config:" + err.Error())
		return
	}

	logger.Info("------------------------------------------------------------------------Gorm数据库连接初始化------------------------------------------------------------------------")
	logger.Info(db.GetDBHandle())
	Config := ConfigManage.GetWebConfig()
	err = db.GetDBHandle().GormDBInit(Config.Database.User, Config.Database.Password, Config.Database.Address, Config.Database.Name, Config.Database.Port)
	if err != nil {
		logger.Error("Failed to init database: %v\n", err)
		return
	}

	logger.Info("------------------------------------------------------------------------Nepcat的api接口初始化------------------------------------------------------------------------")
	NepCatInit.InitAllApis()

	logger.Info("------------------------------------------------------------------------Nepcat的websocket连接初始化------------------------------------------------------------------------")
	Nepcat_ws_init.WebSocketInit(Config.WebsocketInfo.Scheme, Config.WebsocketInfo.Host, Config.WebsocketInfo.Port, Config.WebsocketInfo.Path, Config.WebsocketInfo.Rawquery)
}
