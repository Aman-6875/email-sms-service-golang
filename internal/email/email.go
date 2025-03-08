package email

import (
	"email-sms-service/config"
	"email-sms-service/pkg/logger"
	"fmt"
)

var (
	emailsSentCounter   = 0
	emailsFailedCounter = 0
)

// SendEmail sends an email using the configured provider
func sendEmail(task EmailTask) error {
	provider := getEmailProvider()
	logger.Log.Infof("Using email provider: %s", provider.Name())
	// Send the email
	return provider.SendEmail(task)
}

type Attachment struct {
	FileName string `json:"fileName"` // Name of the file
	Content  []byte `json:"content"`  // File content as bytes
}

type EmailTask struct {
	To           string                 `json:"to"`
	Subject      string                 `json:"subject"`
	Body         string                 `json:"body"`
	Template     string                 `json:"template"`
	TemplateData map[string]interface{} `json:"templateData"`
	Attachments  []*Attachment          `json:"attachments"`
}

func SendEmail(task EmailTask) error {
	const maxRetries = 3
	var err error

	for i := 0; i < maxRetries; i++ {
		err = sendEmail(task)
		if err == nil {
			emailsSentCounter++
			logger.Log.Infof("Email sent successfully to: %s (Total sent: %d)", task.To, emailsSentCounter)
			return nil // Email sent successfully
		}
		logger.Log.Warnf("Attempt %d: Failed to send email: %v", i+1, err)
	}

	emailsFailedCounter++
	logger.Log.Errorf("Failed to send email after %d attempts: %v (Total failed: %d)", maxRetries, err, emailsFailedCounter)

	return fmt.Errorf("failed to send email after %d attempts: %v", maxRetries, err)
}

func getEmailProvider() EmailProvider {
	providerType := config.GetEnv("EMAIL_PROVIDER", "smtp")

	switch providerType {
	case "smtp":
		smtpConfig := SMTPConfig{
			Host:     config.GetEnv("SMTP_HOST", "smtp.example.com"),
			Port:     config.GetEnv("SMTP_PORT", "587"),
			Username: config.GetEnv("SMTP_USERNAME", ""),
			Password: config.GetEnv("SMTP_PASSWORD", ""),
			From:     config.GetEnv("SMTP_FROM", "no-reply@example.com"),
		}
		return NewSMTPProvider(smtpConfig)

	default:
		logger.Log.Fatalf("Unsupported email provider: %s", providerType)
		return nil
	}
}
