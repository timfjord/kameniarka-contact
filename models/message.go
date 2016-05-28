package models

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go"
)

// Message struct
type Message struct {
	Name  string `form:"name" json:"name" binding:"required"`
	Phone string `form:"phone" json:"phone" binding:"required"`
	Email string `form:"email" json:"email"`
	Text  string `form:"text" json:"text" binding:"required"`
}

// Format Message for delivering
func (msg *Message) format() string {
	format := `
<html>
	<body>
		<b>Ім’я:</b> %s<br/>
		<b>Ел.пошта:</b> %s<br/>
		<b>Телефон:</b> %s<br/>
		<b>Повідомлення:</b><br/>
		%s
	</body>
</html>
`
	return fmt.Sprintf(format, msg.Name, msg.Email, msg.Phone, msg.Text)
}

// Deliver via sendgrid
func (msg *Message) Deliver() error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if apiKey == "" {
		return errors.New("SENDGRID_API_KEY is missing")
	}
	email := os.Getenv("CONTACT_EMAIL")
	if email == "" {
		return errors.New("CONTACT_EMAIL is missing")
	}

	sg := sendgrid.NewSendGridClientWithApiKey(apiKey)
	mail := sendgrid.NewMail()
	mail.AddTo(email)
	mail.SetSubject(fmt.Sprintf("Нове повідомлення з сайту від %s", time.Now().Format("02.01.2006")))
	mail.SetFrom("no-reply@kameniarka.cv.ua")
	mail.SetHTML(msg.format())

	return sg.Send(mail)
}
