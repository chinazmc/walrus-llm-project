package route

import (
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"walrus_llm_project/api/controller"
	"walrus_llm_project/api/middleware"
	"walrus_llm_project/log"
)

func Setup(gin *gin.Engine) {
	publicRouter := gin.Group("")
	publicRouter.Use(
		ginzap.RecoveryWithZap(log.Logger.Logger, true),
		cors.Default(),
		middleware.ResponseLogMiddleware(log.Logger),
		middleware.RequestLogMiddleware(log.Logger),
	)
	hc := controller.NewHandleController()
	publicRouter.POST("/v1/handle", hc.HandleQuestion)
}
