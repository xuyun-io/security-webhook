package api

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"security-webhook/api/validating"
	"security-webhook/utils/log"
)

var certFile = "/etc/certs/tls.crt"
var keyFile = "/etc/certs/tls.key"

func InitAPI(e *gin.Engine) {
	e.POST("/validate-privileged-container", validating.PrivilegedContainerCheck)

	err := e.RunTLS(":9443", certFile, keyFile)
	if err != nil {
		log.Logger.Error("failed to start server", zap.Error(err))
	}
}
