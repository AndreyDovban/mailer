package main

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"

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

	_, err = w.Write([]byte("Hello"))
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
