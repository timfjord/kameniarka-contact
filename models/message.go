package models

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

	from := mail.NewEmail("Kameniarka Bot", "no-reply@kameniarka.com")
	subject := fmt.Sprintf("Нове повідомлення з сайту від %s", time.Now().Format("02.01.2006"))
	to := mail.NewEmail("", email)
	message := mail.NewSingleEmail(from, subject, to, "", msg.format())

	client := sendgrid.NewSendClient(apiKey)
	_, err := client.Send(message)

	return err
}
