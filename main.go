package main

import (
	"github.com/gin-gonic/gin"
	"security-webhook/api"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	api.InitAPI(engine)
}
