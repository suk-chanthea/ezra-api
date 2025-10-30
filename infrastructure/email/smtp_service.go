package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailService handles email operations
type EmailService interface {
	SendOTP(to, code, purpose string) error
	SendEmail(to, subject, body string) error
}

type smtpEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// NewSMTPEmailService creates a new SMTP email service
func NewSMTPEmailService(host, port, username, password, from string) EmailService {
	return &smtpEmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// SendOTP sends an OTP code to the specified email
func (s *smtpEmailService) SendOTP(to, code, purpose string) error {
	subject := "Your Verification Code"
	
	var purposeText string
	switch purpose {
	case "email_verification":
		purposeText = "email verification"
	case "password_reset":
		purposeText = "password reset"
	case "login":
		purposeText = "login verification"
	default:
		purposeText = "verification"
	}

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<style>
		body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
		.container { max-width: 600px; margin: 0 auto; padding: 20px; }
		.header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; border-radius: 5px 5px 0 0; }
		.content { background-color: #f9f9f9; padding: 30px; border: 1px solid #ddd; }
		.otp-code { font-size: 32px; font-weight: bold; color: #4CAF50; text-align: center; padding: 20px; background-color: #fff; border: 2px dashed #4CAF50; border-radius: 5px; margin: 20px 0; letter-spacing: 5px; }
		.footer { text-align: center; padding: 20px; font-size: 12px; color: #777; }
		.warning { color: #f44336; font-size: 14px; margin-top: 20px; }
	</style>
</head>
<body>
	<div class="container">
		<div class="header">
			<h1>Verification Code</h1>
		</div>
		<div class="content">
			<h2>Hello!</h2>
			<p>You have requested a verification code for <strong>%s</strong>.</p>
			<p>Please use the following code to complete your verification:</p>
			<div class="otp-code">%s</div>
			<p>This code will expire in <strong>10 minutes</strong>.</p>
			<p class="warning">⚠️ If you did not request this code, please ignore this email and ensure your account is secure.</p>
		</div>
		<div class="footer">
			<p>This is an automated email. Please do not reply.</p>
			<p>&copy; 2025 Ezra. All rights reserved.</p>
		</div>
	</div>
</body>
</html>
	`, purposeText, code)

	return s.SendEmail(to, subject, body)
}

// SendEmail sends a generic email
func (s *smtpEmailService) SendEmail(to, subject, body string) error {
	// Setup authentication
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	// Compose message with HTML
	headers := make(map[string]string)
	headers["From"] = s.from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"utf-8\""

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP server
	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	
	// Send the email
	err := smtp.SendMail(
		addr,
		auth,
		s.from,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

// ValidateEmail performs basic email validation
func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

