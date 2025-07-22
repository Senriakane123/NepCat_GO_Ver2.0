package Init

import (
	"ResourceKeeper/Gin_Init/Verify"
	"github.com/gin-gonic/gin"
)

//var RestInit sync.Map

func RestfullInit(router *gin.Engine) {
	MapStrFuncRouter := map[string][]func(*gin.RouterGroup){}

	for groupstr, funcs := range MapStrFuncRouter {
		for _, f := range funcs {
			var Router *gin.RouterGroup
			if groupstr == "/cdp/v1/login" {
				Router = router.Group(groupstr)
			} else {
				// 其他接口带上中间件
				Router = router.Group(groupstr, Verify.CDPVerify())
			}

			f(Router)
		}
	}
}
