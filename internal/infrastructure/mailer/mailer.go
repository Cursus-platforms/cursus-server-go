package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"strconv"
)

type EmailData struct {
	Name    string
	CodeOTP string
}

const templatePath = "templates/otp_verification.html"

type Mailer interface {
	SendOTP(recipientEmail, recipientName, otp string) error
}

type SMTPMailer struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewSMTPMailer() Mailer {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if port == 0 {
		port = 587
	}

	return &SMTPMailer{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("FROM_EMAIL"),
		Password: os.Getenv("FROM_EMAIL_PASSWORD"),
	}
}

func (m *SMTPMailer) SendOTP(recipientEmail, recipientName, otp string) error {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("mailer: lỗi parse template: %w", err)
	}

	data := EmailData{Name: recipientName, CodeOTP: otp}
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("mailer: lỗi execute template: %w", err)
	}

	var msg bytes.Buffer

	fromAddress := os.Getenv("MAIL_FROM_ADDRESS")
	if fromAddress == "" {
		fromAddress = m.Username
	}

	msg.WriteString(fmt.Sprintf("From: %s\r\n", fromAddress))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", recipientEmail))
	msg.WriteString("Subject: Mã xác minh tài khoản của bạn (OTP)\r\n")
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	msg.WriteString("\r\n") // Kết thúc Headers

	msg.Write(body.Bytes())

	smtpAddr := fmt.Sprintf("%s:%d", m.Host, m.Port)
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	err = smtp.SendMail(
		smtpAddr,
		auth,
		fromAddress,
		[]string{recipientEmail},
		msg.Bytes(),
	)

	if err != nil {
		return fmt.Errorf("mailer: lỗi gửi email qua net/smtp: %w", err)
	}

	return nil
}
