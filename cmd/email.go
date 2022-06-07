/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

const MAILDATE string = "2006/01/02-15:04:05"

func SendLogEmail(subject string, body string, attachments ...string) error {
	server, err := GetEnvVar("ERRWRAPPER_MAIL_SERVER")
	if err != nil {
		return err
	}

	portString, err := GetEnvVar("ERRWRAPPER_MAIL_PORT")
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(portString)
	if err != nil {
		return err
	}

	from, err := GetEnvVar("ERRWRAPPER_MAIL_FROM")
	if err != nil {
		return err
	}

	to, err := GetEnvVar("ERRWRAPPER_MAIL_TO")
	if err != nil {
		return err
	}

	user, err := GetEnvVar("ERRWRAPPER_MAIL_USER")
	if err != nil {
		return err
	}

	pass, err := GetEnvVar("ERRWRAPPER_MAIL_PASS")
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	for _, attachment := range attachments {
		filestat, err := os.Stat(attachment)
		if err != nil {
			return fmt.Errorf("file %q to be attached does not exist", attachment)
		}
		if filestat.Size() > 0 {
			m.Attach(attachment)
		}
	}

	d := gomail.NewDialer(server, port, user, pass)

	err = d.DialAndSend(m)
	if err != nil {
		return errors.New("failed to send email")
	}

	return nil
}
