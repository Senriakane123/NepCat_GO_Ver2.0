package Init

import (
	"ResourceKeeper/ConfigManage"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Gin_Init() {
	fmt.Println("------------------------------------------------------------------------Gin初始化------------------------------------------------------------------------")
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https://yourdomain.com"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	RestfullInit(router)

	// 配置 HTTP 或 HTTPS 服务
	if ConfigManage.GetWebConfig().Server.HTTPEnabled {
		// 启动 HTTP 服务
		fmt.Println("------------------------------------------------------------------------启动http服务------------------------------------------------------------------------")

		if err := router.Run(fmt.Sprintf(":%d", ConfigManage.GetWebConfig().Server.HTTPPort)); err != nil {
			fmt.Printf("Failed to start HTTP server: %v\n", err)
		}

	}

	if ConfigManage.GetWebConfig().Server.HTTPSEnabled {
		// 启动 HTTPS 服务
		fmt.Println("------------------------------------------------------------------------启动https服务------------------------------------------------------------------------")

		if err := router.RunTLS(fmt.Sprintf(":%d", ConfigManage.GetWebConfig().Server.HTTPSPort), "server.crt", "server.key"); err != nil {
			fmt.Printf("Failed to start HTTPS server: %v\n", err)
		}

	}
}
