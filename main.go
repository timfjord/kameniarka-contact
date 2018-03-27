package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/haisum/recaptcha"
	"github.com/itsjamie/gin-cors"
	"github.com/timsly/kameniarka-contact/models"
)

func main() {
	r := gin.New()

	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	r.GET("/contact", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
	})

	r.POST("/contact", func(c *gin.Context) {
		re := recaptcha.R{ Secret: os.Getenv("RECAPTCHA_SECRET") }
		if re.VerifyResponse(c.Request.FormValue("g-recaptcha-response")) {
			var msg models.Message
			if c.Bind(&msg) == nil {
				if e := msg.Deliver(); e == nil {
					c.Writer.WriteHeader(http.StatusCreated)
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
				}
			} else {
				c.Writer.WriteHeader(http.StatusBadRequest)
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": re.LastError()})
		}
	})

	r.Run()
}
