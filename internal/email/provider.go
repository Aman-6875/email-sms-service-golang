package email

import (
	"email-sms-service/pkg/logger"
	"fmt"
	"io"

	"gopkg.in/gomail.v2"
)

// EmailProvider defines the interface for sending emails
type EmailProvider interface {
	Name() string
	SendEmail(task EmailTask) error
}

// SMTPConfig holds the configuration for the SMTP provider
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// SMTPProvider implements the EmailProvider interface using SMTP
type SMTPProvider struct {
	config SMTPConfig
}

// Name returns the name of the provider
func (p *SMTPProvider) Name() string {
	return "SMTP"
}

// NewSMTPProvider creates a new SMTPProvider instance
func NewSMTPProvider(config SMTPConfig) *SMTPProvider {
	return &SMTPProvider{
		config: config,
	}
}

// SendEmail sends an email using the SMTP provider
func (p *SMTPProvider) SendEmail(task EmailTask) error {
	const maxRetries = 3
	var err error

	for i := 0; i < maxRetries; i++ {
		err = p.sendEmail(task)
		if err == nil {
			return nil // Email sent successfully
		}

		logger.Log.Warnf("Attempt %d: Failed to send email: %v", i+1, err)
	}

	return fmt.Errorf("failed to send email after %d attempts: %v", maxRetries, err)
}

// sendEmail contains the actual email sending logic
func (p *SMTPProvider) sendEmail(task EmailTask) error {
	// Parse the HTML template if provided
	var htmlBody string
	if task.Template != "" {
		var err error
		htmlBody, err = ParseTemplate(task.Template, task.TemplateData)
		if err != nil {
			return fmt.Errorf("failed to parse HTML template: %v", err)
		}
	}

	// Create a new email message
	m := gomail.NewMessage()
	m.SetHeader("From", p.config.From)
	m.SetHeader("To", task.To)
	m.SetHeader("Subject", task.Subject)

	// Add plain text body
	if task.Body != "" {
		m.SetBody("text/plain", task.Body)
	}

	// Add HTML body
	if htmlBody != "" {
		m.SetBody("text/html", htmlBody)
	}

	// Add attachments
	for _, attachment := range task.Attachments {
		m.Attach(attachment.FileName, gomail.SetCopyFunc(func(w io.Writer) error {
			_, err := w.Write(attachment.Content)
			return err
		}))
	}

	// Create a new dialer
	d := gomail.NewDialer(p.config.Host, 465, p.config.Username, p.config.Password)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	logger.Log.Infof("Email sent successfully to: %s", task.To)
	return nil
}