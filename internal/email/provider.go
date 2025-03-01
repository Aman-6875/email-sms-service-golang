package email

import (
	"crypto/tls"
	"email-sms-service/pkg/logger"
	"fmt"
	"net/smtp"
)

type EmailProvider interface {
    Name() string
	SendEmail(to, subject, body string) error
}


type SMTPConfig struct {
	Host string
	Port string
	Username string
    Password string
    From     string
}

type SMTPProvider struct {
    config SMTPConfig
}

func (p *SMTPProvider) Name() string {
    return "SMTP"
}
func NewSMTPProvider (config SMTPConfig) *SMTPProvider {
	return &SMTPProvider{
		config: config,
	}
}

func (p *SMTPProvider) SendEmail(to, subject, body string) error {
	auth := smtp.PlainAuth("", p.config.Username, p.config.Password, p.config.Host)

	client, err := smtp.Dial(fmt.Sprintf("%s:%s", p.config.Host, p.config.Port))
    if err != nil {
        return fmt.Errorf("failed to connect to SMTP server: %v", err)
    }
    defer client.Quit()

    // Start TLS
    if err := client.StartTLS(&tls.Config{
        InsecureSkipVerify: true, // Only for testing; use proper certificates in production
        ServerName:         p.config.Host,
    }); err != nil {
        return fmt.Errorf("failed to start TLS: %v", err)
    }
	
	if err := client.Auth(auth); err != nil {
        return fmt.Errorf("failed to authenticate: %v", err)
    }

	if err := client.Mail(p.config.From); err != nil {
        return fmt.Errorf("failed to set sender: %v", err)
    }

	if err := client.Rcpt(to); err != nil {
        return fmt.Errorf("failed to set recipient: %v", err)
    }

	wc, err := client.Data()
    if err != nil {
        return fmt.Errorf("failed to prepare email body: %v", err)
    }
    defer wc.Close()

	emailBody := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	if _, err := wc.Write([]byte(emailBody)); err != nil {
        return fmt.Errorf("failed to write email body: %v", err)
    }

    logger.Log.Infof("Email sent successfully to: %s", to)
    return nil
}
