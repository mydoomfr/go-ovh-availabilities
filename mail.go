package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
	"os/exec"
)

type Request struct {
	from    string
	to      []string
	subject string
	body    string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

func NewRequest(to []string, subject string) *Request {
	return &Request{
		to:      to,
		subject: subject,
	}
}

func (r *Request) parseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}
	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}
	r.body = buffer.String()
	return nil
}

func (r *Request) sendMailWithSmtpAuth() bool {
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	SMTP := fmt.Sprintf("%s:%d", config.Mail.SMTP.Server, config.Mail.SMTP.Port)
	AUTH := smtp.PlainAuth(
		"", config.Mail.SMTP.Username, config.Mail.SMTP.Password, config.Mail.SMTP.Server)
	if err := smtp.SendMail(SMTP, AUTH, config.Mail.From, r.to, []byte(body)); err != nil {
		return false
	}
	return true
}

func (r *Request) sendMail() {
	body := "To: " + r.to[0] + "\r\nSubject: " + r.subject + "\r\n" + MIME + "\r\n" + r.body
	sendmail := exec.Command(config.Mail.Sendmail.Bin, "-f", config.Mail.From, config.Mail.To)
	stdin, err := sendmail.StdinPipe()
	if err != nil {
		panic(err)
	}
	stdout, err := sendmail.StdoutPipe()
	if err != nil {
		panic(err)
	}
	sendmail.Start()
	stdin.Write([]byte(body))
	stdin.Close()
	sentBytes, _ := ioutil.ReadAll(stdout)
	sendmail.Wait()
	log.Printf("Mail has been sent to %s. %s\n", r.to, string(sentBytes))
}

func (r *Request) Send(templateName string, items interface{}) {
	err := r.parseTemplate(templateName, items)
	if err != nil {
		log.Fatal(err)
	}
	// SMTP Auth
	if config.Mail.SMTP.Active {
		if ok := r.sendMailWithSmtpAuth(); ok {
			log.Printf("Mail has been sent to %s\n", r.to)
		} else {
			log.Fatalf("Failed to send the mail to %s\n", r.to)
		}
	} else if config.Mail.Sendmail.Active {
		r.sendMail()
	} else {
		log.Printf("Mail support disabled.\n")
	}
}
