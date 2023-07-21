package utils

import (
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
)

type Mail struct {
	Sender  string
	To      []string
	Subject string
	Body    string
}

func (mail *Mail) BuildMessage() string {
	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\r\n"
	msg += fmt.Sprintf("From: %s\r\n", mail.Sender)
	msg += fmt.Sprintf("To: %s\r\n", strings.Join(mail.To, ";"))
	msg += fmt.Sprintf("Subject: %s\r\n", mail.Subject)
	msg += fmt.Sprintf("\r\n%s\r\n", mail.Body)
	return msg
}

type AuthTokenJwtClaim struct {
	Email string
	ID    uint
	jwt.StandardClaims
}
