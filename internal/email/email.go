package email

import (
    "email-sms-service/config"
    "email-sms-service/pkg/logger"
)

// SendEmail sends an email using the configured provider
func SendEmail(to, subject, body string) error {
    // Load SMTP configuration from environment variables
    smtpConfig := SMTPConfig{
        Host:     config.GetEnv("SMTP_HOST", "smtp.example.com"),
        Port:     config.GetEnv("SMTP_PORT", "587"),
        Username: config.GetEnv("SMTP_USERNAME", ""),
        Password: config.GetEnv("SMTP_PASSWORD", ""),
        From:     config.GetEnv("SMTP_FROM", "no-reply@example.com"),
    }

    // Create a new SMTP provider
    provider := NewSMTPProvider(smtpConfig)

    // Send the email
    err := provider.SendEmail(to, subject, body)
    if err != nil {
        logger.Log.Errorf("Failed to send email: %v", err)
        return err
    }

    return nil
}