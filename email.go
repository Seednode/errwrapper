/*
Copyright © 2026 Seednode <seednode@seedno.de>
*/

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

const MAILDATE string = "2006/01/02-15:04:05"

func SendLogEmail(subject string, body string, attachments ...string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", mailFrom)
	m.SetHeader("To", mailTo)
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

	port, err := strconv.Atoi(mailPort)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(mailServer, port, mailUser, mailPass)

	err = d.DialAndSend(m)
	if err != nil {
		return errors.New("failed to send email")
	}

	return nil
}
