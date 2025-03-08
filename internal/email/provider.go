package email

import (
	"bytes"
	"email-sms-service/pkg/logger"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
)

type EmailProvider interface {
	Name() string
	SendEmail(task EmailTask) error
}

type SMTPConfig struct {
	Host     string
	Port     string
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
func NewSMTPProvider(config SMTPConfig) *SMTPProvider {
	return &SMTPProvider{
		config: config,
	}
}

func (p *SMTPProvider) SendEmail(task EmailTask) error {
	var htmlBody string

	if task.Template != "" {
		var err error
		htmlBody, err = ParseTemplate(task.Template, task.TemplateData)

		if err != nil {
			return fmt.Errorf("failed to parse HTML template: %v", err)
		}
	}
	auth := smtp.PlainAuth("", p.config.Username, p.config.Password, p.config.Host)

	var msg bytes.Buffer

	writer := multipart.NewWriter(&msg)

	headers := map[string]string{
		"From":         p.config.From,
		"To":           task.To,
		"Subject":      task.Subject,
		"MIME-Version": "1.0",
		"Content-Type": "multipart/mixed; boundary=" + writer.Boundary(),
	}

	for key, value := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
	}
	msg.WriteString("\r\n")

	if task.Body != "" {
		part, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Type": []string{"text/plain; charset=UTF-8"},
		})
		if err != nil {
			return fmt.Errorf("failed to create plain text part: %v", err)
		}
		part.Write([]byte(htmlBody))
	}

    if htmlBody != "" {
		part, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Type": []string{"text/html; charset=UTF-8"},
		})
		if err != nil {
			return fmt.Errorf("failed to create HTML part: %v", err)
		}
		part.Write([]byte(htmlBody))
	}

	for _, attachment := range task.Attachments {
		part, err := writer.CreatePart(textproto.MIMEHeader{
			"Content-Type":        []string{"application/octet-stream"},
			"Content-Disposition": []string{fmt.Sprintf("attachment; filename=\"%s\"", attachment.FileName)},
		})
		if err != nil {
			return fmt.Errorf("failed to create attachment part: %v", err)
		}
		part.Write(attachment.Content)
	}
	writer.Close()

	err := smtp.SendMail(
		fmt.Sprintf("%s:%s", p.config.Host, p.config.Port),
		auth,
		p.config.From,
		[]string{task.To},
		msg.Bytes(),
	)

	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	logger.Log.Infof("Email sent successfully to: %s", task.To)
	return nil
	// client, err := smtp.Dial(fmt.Sprintf("%s:%s", p.config.Host, p.config.Port))
	// if err != nil {
	//     return fmt.Errorf("failed to connect to SMTP server: %v", err)
	// }
	// defer client.Quit()

	// // Start TLS
	// if err := client.StartTLS(&tls.Config{
	//     InsecureSkipVerify: true, // Only for testing; use proper certificates in production
	//     ServerName:         p.config.Host,
	// }); err != nil {
	//     return fmt.Errorf("failed to start TLS: %v", err)
	// }

	// if err := client.Auth(auth); err != nil {
	//     return fmt.Errorf("failed to authenticate: %v", err)
	// }

	// if err := client.Mail(p.config.From); err != nil {
	//     return fmt.Errorf("failed to set sender: %v", err)
	// }

	// if err := client.Rcpt(to); err != nil {
	//     return fmt.Errorf("failed to set recipient: %v", err)
	// }

	// wc, err := client.Data()
	// if err != nil {
	//     return fmt.Errorf("failed to prepare email body: %v", err)
	// }
	// defer wc.Close()

	// emailBody := fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body)

	// if _, err := wc.Write([]byte(emailBody)); err != nil {
	//     return fmt.Errorf("failed to write email body: %v", err)
	// }

	// logger.Log.Infof("Email sent successfully to: %s", to)
	// return nil
}
