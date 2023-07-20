package config

import (
	"net/smtp"
	"crypto/tls"
	"fmt"
)


func SendMail(rcpt, message string) bool {
	mailConfig := Config.Mail
	fmt.Println("Mail config:", mailConfig)

	var auth smtp.Auth = nil
	if mailConfig.UseAuth {
		auth = smtp.PlainAuth("", mailConfig.Username, mailConfig.Password, mailConfig.Hostname)
		// auth := smtp.CRAMMD5Auth(mailConfig.Username, mailConfig.Password)
	}

	messageBytes := []byte(message)

	rcpts := []string{rcpt}

	// SEND EMAIL THE EASY WAY
	if err := smtp.SendMail(mailConfig.Hostname, auth, mailConfig.From, rcpts, messageBytes); err == nil {
		return true
	} else {
		fmt.Println("Error sending mail:", err)
	}


	// SEND EMAIL THE HARD WAY
	// connect
	client, err := smtp.Dial(mailConfig.Hostname)
	if err != nil {
		fmt.Println("Error sending mail:", err)
		return false
	}
	defer client.Quit()
	
	if hasStartTLS, _ := client.Extension("STARTTLS"); mailConfig.UseStartTLS && hasStartTLS {
		tlsConfig := &tls.Config{ServerName: mailConfig.Hostname, InsecureSkipVerify: true} 
		client.StartTLS(tlsConfig)
	}

	if mailConfig.UseAuth {
		err = client.Auth(auth)
		if err != nil {
			fmt.Println("Error sending mail:", err)
			return false
		}
	}

	// set recipient
	client.Mail(mailConfig.From)
	client.Rcpt(rcpt)

	// write the body
	writer, err := client.Data()
	if err != nil {
		fmt.Println("Error sending mail:", err)
		return false
	}
	_, err = fmt.Fprintf(writer, message)
	if err != nil {
		fmt.Println("Error sending mail:", err)
		return false
	}
	err = writer.Close()
	if err != nil {
		fmt.Println("Error sending mail:", err)
		return false
	}

	// Send the QUIT command and close the connection.
	err = client.Quit()
	if err != nil {
		fmt.Println("Error sending mail:", err)
		return false
	}

	return true
}