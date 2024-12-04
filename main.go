package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

func main() {

	fmt.Print("\033[H\033[2J")

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file? using default config")
	}

	sender := os.Getenv("LOGIN")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := 587
	address := fmt.Sprintf("%s:%d", host, port)
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	auth := smtp.PlainAuth("", sender, password, host)

	var body bytes.Buffer
	t, err := template.ParseFiles("./test.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	t.Execute(&body, struct{ Name string }{Name: "Andrey"})

	from := mail.Address{Name: "Name", Address: sender}
	to := mail.Address{Name: "", Address: sender}
	subj := "This is the email subject"
	mimeVersion := "1.0;"
	contentType := "text/html; charset=\"UTF-8\";"

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj
	headers["MIME-version"] = mimeVersion
	headers["Content-Type"] = contentType

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body.String()

	fmt.Println(message)

	c, err := smtp.Dial(address)
	if err != nil {
		fmt.Println(err.Error())
	}

	if err = c.StartTLS(tlsconfig); err != nil {
		fmt.Println(err.Error())
		c.Close()
	}

	if err = c.Auth(auth); err != nil {
		fmt.Println(err.Error())
	}

	if err = c.Mail(sender); err != nil {
		fmt.Println(err.Error())
		c.Close()
	}

	if err = c.Rcpt(sender); err != nil {
		fmt.Println(err.Error())
		c.Close()
	}

	w, err := c.Data()

	if err != nil {
		fmt.Println(err.Error())
		c.Close()
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		fmt.Println(err.Error())
		c.Close()
	}

	err = w.Close()

	if err != nil {
		fmt.Println(err.Error())
		c.Close()
	}

	fmt.Println("Send mail success!")
	c.Quit()
}
