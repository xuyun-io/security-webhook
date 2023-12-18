package main

import (
	"github.com/gin-gonic/gin"
	"security-webhook/api"
	"security-webhook/utils/log"
)

func main() {
	defer log.Logger.Sync()

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	api.InitAPI(engine)
}
