package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"maildefender/notifier/internal/api/handlers"
	"maildefender/notifier/internal/client"
)

var engine *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
	engine = gin.New()
}

var additionalVariables map[string]any = map[string]any{}

func SetSmtpClient(client client.Client) {
	additionalVariables["smtpClient"] = client
}

func RegisterMiddlewares() {
	engine.Use(logger())
}

func logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		for key, value := range additionalVariables {
			c.Set(key, value)
		}

		t := time.Now()

		// before request

		c.Next()

		// after request
		latency := time.Since(t)

		// access the status we are sending
		status := c.Writer.Status()

		logger := logrus.WithFields(logrus.Fields{
			"status":  status,
			"latency": latency,
			"method":  c.Request.Method,
		})

		logBody := fmt.Sprintf("Handling %s", c.Request.URL.Path)

		switch {
		case status < 300:
			logger.Info(logBody)
			break
		case status >= 300 && status < 400:
			logger.Warn(logBody)
			break
		case status >= 400:
			logger.Error(logBody)
		}
	}
}

func RegisterHandlers() {
	v1 := engine.Group("/v1/notifier")

	v1.POST("/email", handlers.SendMail)
}

func Run() error {
	return engine.Run(":8080")
}
