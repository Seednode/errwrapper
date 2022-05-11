/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

const MAILDATE string = "2006/01/02-15:04:05"

func SendLogEmail(subject string, body string, attachments ...string) {
	server := GetEnvVar("ERRWRAPPER_MAIL_SERVER")
	port, _ := strconv.Atoi(GetEnvVar("ERRWRAPPER_MAIL_PORT"))
	from := GetEnvVar("ERRWRAPPER_MAIL_FROM")
	to := GetEnvVar("ERRWRAPPER_MAIL_TO")
	user := GetEnvVar("ERRWRAPPER_MAIL_USER")
	pass := GetEnvVar("ERRWRAPPER_MAIL_PASS")

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	for _, attachment := range attachments {
		filestat, err := os.Stat(attachment)
		if err != nil {
			panic(err)
		}
		if filestat.Size() > 0 {
			m.Attach(attachment)
		}
	}

	d := gomail.NewDialer(server, port, user, pass)

	err := d.DialAndSend(m)
	if err != nil {
		panic(err)
	}
}
