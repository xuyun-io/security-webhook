package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"security-webhook/utils/log"
)

var certFile = "/etc/certs/tls.crt"
var keyFile = "/etc/certs/tls.key"

func InitAPI(e *gin.Engine) {
	//e.POST("/validate-privileged-container", validating.PrivilegedContainerCheck)
	//
	//err := e.RunTLS(":9443", certFile, keyFile)
	////err := e.Run(":9443")
	//if err != nil {
	//	log.Logger.Error("failed to start server", zap.Error(err))
	//}

	e.GET("/proxy", gin.HandlerFunc(func(context *gin.Context) {
		err := context.Request.ParseForm()
		if err != nil {
			log.Logger.Error("ParseForm" + err.Error())
			return
		}

		target := context.Request.FormValue("https")
		fmt.Printf("---->> https://%s", target)
		tt := fmt.Sprintf("https://%s", target)
		rsp, err := http.DefaultClient.Get(tt)
		if err != nil {
			log.Logger.Error("Get" + err.Error())
			return
		}

		data, err := io.ReadAll(rsp.Body)
		if err != nil {
			log.Logger.Error("ReadAll" + err.Error())
			return
		}
		context.String(200, "%s", string(data))
	}))

	err := e.Run(":80")
	if err != nil {
		log.Logger.Error("failed to start server", zap.Error(err))
	}
}
