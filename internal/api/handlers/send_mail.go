package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"maildefender/notifier/internal/client"
	"maildefender/notifier/internal/configuration"
	"maildefender/notifier/internal/models"
)

type sendMailIn struct {
	To          []string `json:"to"`
	Subject     string   `json:"subject"`
	Body        string   `json:"body"`
	ReplyTo     string   `json:"replyTo"`
	ThreadTopic string   `json:"threadTopic"`
}

func SendMail(c *gin.Context) {
	smtpClient := c.MustGet("smtpClient").(client.Client)

	var in sendMailIn

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := smtpClient.Send(models.MailStructure{
		To:          in.To,
		Subject:     in.Subject,
		ReplyTo:     in.ReplyTo,
		ThreadTopic: in.ThreadTopic,
		Body:        template.HTML(in.Body),
		From:        configuration.SmtpConfiguration().Authentication.Username,
		Date:        time.Now(),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, nil)
}
