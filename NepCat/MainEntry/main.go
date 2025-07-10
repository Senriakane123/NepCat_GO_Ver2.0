package main

import (
	ConfigManage "NepCat_GO/ConfigModule"
	"NepCat_GO/NepCatInit"
	"NepCat_GO/NepCatInit/Message_Handle"
	"NepCat_GO/NepCatInit/MsgProcess"
	"NepCat_GO/NepCatInit/Nepcat_ws_init"
	db "github.com/gopublic/GormModule/DBControl/DatabaseControl"
	"github.com/jander/golog/logger"
	"sync"
)

func main() {
	rotatingHandler := logger.NewFileHandler(ConfigManage.LogFilePath)
	logger.SetHandlers(logger.Console, rotatingHandler)
	logger.Info("------------------------------------------------------------------------系统日志初始化------------------------------------------------------------------------")

	logger.Info("------------------------------------------------------------------------config初始化------------------------------------------------------------------------")
	err := ConfigManage.ConfigInit(ConfigManage.ConfigFilePath) // 这里假设配置文件名为 config.yaml
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
	MsgProcess.RespApiInit()
	NepCatInit.InitAllApis()
	// WebSocket + 消息处理线程

	logger.Info("------------------------------------------------------------------------开启消息处理线程------------------------------------------------------------------------")
	Nepcat_ws_init.WebChannelInit()
	go Message_Handle.MessageHandler()
	Nepcat_ws_init.WebSocketInit(Config.WebsocketInfo.Scheme, Config.WebsocketInfo.Host, Config.WebsocketInfo.Port, Config.WebsocketInfo.Path, Config.WebsocketInfo.Rawquery)

	// 阻塞主线程
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
