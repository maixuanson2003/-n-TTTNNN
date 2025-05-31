package Config

import (
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func newEmailConfig() EmailConfig {
	Env := GetEnvConfig()
	port, err := strconv.Atoi(Env.SmtpPort())
	if err != nil {
		log.Print(err)
		return EmailConfig{}
	}
	return EmailConfig{
		Host:     Env.SmtpHost(),
		Port:     port,
		Username: Env.SmtpUser(),
		Password: Env.SmtpPassword(),
		From:     Env.FromEmail(),
	}
}
func SendEmail(to, subject, body string) error {
	config := newEmailConfig()
	email := gomail.NewMessage()
	email.SetHeader("From", config.From)
	email.SetHeader("To", to)
	email.SetHeader("Subject", subject)
	email.SetBody("text/plain", body)
	dialer := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	err := dialer.DialAndSend(email)
	if err != nil {
		log.Print(err)
		return err
	}
	log.Print("email gui thanh cong")
	return nil
}
