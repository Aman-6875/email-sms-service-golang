package main

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

func main() {
	// Create a new message
	m := gomail.NewMessage()
	m.SetHeader("From", "gomail@softstation.xyz")
	m.SetHeader("To", "recipient@example.com")
	m.SetHeader("Subject", "Test Email")
	m.SetBody("text/plain", "This is a test email.")

	// Create a new dialer with SSL/TLS settings
	d := gomail.NewDialer("mail.softstation.xyz", 465, "gomail@softstation.xyz", "gomail@softstation.xyz")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true} // Only for testing; use proper certificates in production

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		fmt.Printf("Failed to send email: %v\n", err)
	} else {
		fmt.Println("Email sent successfully!")
	}
}
