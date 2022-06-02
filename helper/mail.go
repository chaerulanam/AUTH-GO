package helper

import (
	"auth/config"
	"log"

	"gopkg.in/gomail.v2"
)

func KirimAktifasi(sub string, msg string, addr string) {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", config.SENDER_NAME)
	mailer.SetHeader("To", addr)
	mailer.SetAddressHeader("Cc", "anakkendali01@gmail.com", "Tra Lala La")
	mailer.SetHeader("Subject", sub)
	mailer.SetBody("text/html", msg)

	dialer := gomail.NewDialer(
		config.SMTP_HOST,
		config.SMTP_PORT,
		config.AUTH_EMAIL,
		config.AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Println(err.Error())
	}

	log.Println("Mail sent!")
}
