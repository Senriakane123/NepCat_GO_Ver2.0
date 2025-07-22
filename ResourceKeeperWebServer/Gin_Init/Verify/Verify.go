package Verify

import (
	"ResourceKeeper/Gin_Init/Token_Manage"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CDPVerify() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.GetHeader("Token")
		//update := true

		ok, tk := Token_Manage.Instance().Verify_Token(token)

		if ok && tk.Get_Address() == context.ClientIP() {
			context.Next()
		} else {
			if !ok {
				//VRTSTokenMgr.Instance().DumpToken()
				fmt.Println("Token ", token, " not exist!")
			} else {
				fmt.Println("Token ", token, " address mismatch[", tk.Get_Address(), ":", context.ClientIP(), "]")
			}

			context.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"message": "invalid token or IP address mismatch",
			})
			context.Abort()
		}
	}
}
