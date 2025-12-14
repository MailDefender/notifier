package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"maildefender/notifier/internal/api/handlers"
	"maildefender/notifier/internal/client"
)

var engine *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	engine = gin.New()
	engine.Use(gin.LoggerWithWriter(logrus.New().Writer()))
}

var additionalVariables map[string]any = map[string]any{}

func SetSmtpClient(client client.Client) {
	additionalVariables["smtpClient"] = client
}

func RegisterMiddlewares() {
	engine.Use(func(ctx *gin.Context) {
		for key, value := range additionalVariables {
			ctx.Set(key, value)
		}

		ctx.Next()
	})
}

func RegisterHandlers() {
	v1 := engine.Group("/v1/notifier")

	v1.POST("/email", handlers.SendMail)
}

func Run() error {
	return engine.Run(":8080")
}
