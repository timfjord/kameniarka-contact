package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/timsly/gin-recaptcha"
	"github.com/timsly/kameniarka-contact/models"
)

func main() {
	r := gin.New()

	log.Println("start")

	log.Println("cors use")
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	log.Println("rc use")
	r.Use(recaptcha.Middleware(recaptcha.Config{
		Secret: os.Getenv("RECAPTCHA_SECRET"),
	}))

	r.GET("/contact", func(c *gin.Context) {
		log.Println("ping")
		c.Writer.WriteHeader(http.StatusOK)
	})

	r.POST("/contact", func(c *gin.Context) {
		log.Println("post")
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
	})

	r.Run()
}
