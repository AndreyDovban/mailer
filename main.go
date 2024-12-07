package main

import (
	"bytes"
	"fmt"
	"log"
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

	err = smtp.SendMail(host+":25", auth, sender, []string{sender}, []byte(message))
	if err != nil {
		log.Fatal(err)
	}
}
