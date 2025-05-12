package mail

import (
	"fmt"
	"net/smtp"
)

// SendOTPEmail - отправляет OTP на email
func SendOTPEmail(email, otp string) error {
	from := "maks.karzhaubaev@mail.ru"
	password := "bTnDaKD8ahncvZDK7JFQ"
	to := []string{email}

	// SMTP сервер
	smtpHost := "smtp.mail.ru"
	smtpPort := "587"

	// Сообщение
	subject := "Subject: Your OTP Code"
	body := fmt.Sprintf("Your OTP code is: %s", otp)
	message := []byte(subject + "\n\n" + body)

	// Авторизация
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Отправка почты
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}
	return nil
}
