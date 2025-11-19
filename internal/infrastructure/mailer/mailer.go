package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// Define the template path relative to the project root
const templatePath = "templates/otp_verification.html"

// Mailer interface must be updated to accept the name
type Mailer interface {
	SendOTP(recipientEmail, recipientName, otp string) error
}

// SMTPMailer là triển khai sử dụng gomail
type SMTPMailer struct {
	Host     string
	Port     int
	Username string
	Password string
}

// NewSMTPMailer khởi tạo Mailer từ biến môi trường
func NewSMTPMailer() Mailer {
	// (Bạn PHẢI thêm các biến này vào .env: SMTP_HOST, SMTP_PORT, SMTP_USER, SMTP_PASS)
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	return &SMTPMailer{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASS"),
	}
}

type EmailData struct {
	Name    string
	CodeOTP string
}

func (m *SMTPMailer) SendOTP(recipientEmail, recipientName, otp string) error {

	// 1. Load and parse the HTML template
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("lỗi parse template: %w", err)
	}

	// 2. Prepare dynamic data
	data := EmailData{
		Name:    recipientName,
		CodeOTP: otp,
	}

	// 3. Execute the template into a buffer
	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return fmt.Errorf("lỗi execute template: %w", err)
	}

	// 4. Create the Gomail message
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Username)
	msg.SetHeader("To", recipientEmail)
	msg.SetHeader("Subject", "Mã xác minh tài khoản của bạn (OTP)")

	// Set the body as HTML
	msg.SetBody("text/html", body.String()) // <-- Send the rendered HTML

	// 5. Send the email (Logic remains the same)
	d := gomail.NewDialer(m.Host, m.Port, m.Username, m.Password)

	if err := d.DialAndSend(msg); err != nil {
		return fmt.Errorf("lỗi gửi email: %w", err)
	}

	return nil
}
