package smtp

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/smtp"
)

var (
	data, _ = godotenv.Read(".env")
	sender  = data["EMAIL_SEND_LOGIN"]
	pass    = data["EMAIL_SEND_PASS"]
	auth    = smtp.PlainAuth("", sender, pass, "smtp.gmail.com")
)

func SendPingErrorMail(url, email string) {
	to := []string{email}
	s := fmt.Sprintf("Subject: Error with ping server\n\n"+
		"Hello, we encountered a problem when sending a request to the url \"%s.\""+
		"The server does not seem to be responding", url)
	msg := []byte(s)
	if err := smtp.SendMail("smtp.gmail.com:587", auth, sender, to, msg); err != nil {
		log.Fatal(err)
	}
}
