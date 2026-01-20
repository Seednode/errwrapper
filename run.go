/*
Copyright © 2025 Seednode <seednode@seedno.de>
*/

package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

func RunCommand(arguments []string) error {
	defer HandleExit()

	var err error

	timeZone := os.Getenv("TZ")
	if timeZone != "" {
		time.Local, err = time.LoadLocation(timeZone)
		if err != nil {
			return err
		}
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, os.Interrupt)

	hostName, err := os.Hostname()
	if err != nil {
		return errors.New("unable to retrieve hostname")
	}

	startTime := time.Now()
	stdOutFile, stdErrFile, exitCode, err := LogCommand(arguments)
	if err != nil {
		return err
	}
	stopTime := time.Now()

	command := strings.Join(arguments, " ")

	var wg sync.WaitGroup

	if database {

		wg.Go(func() {

			if databaseType != "cockroachdb" && databaseType != "postgresql" {
				fmt.Println("invalid database type specified")
				return
			}

			databaseURL, err := GetDatabaseURL()
			if err != nil {
				fmt.Println(err)
				return
			}

			sqlStatement, err := CreateSQLStatement(startTime, stopTime, hostName, command, exitCode)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = WriteToDatabase(databaseURL, sqlStatement)
			if err != nil {
				fmt.Println(err)
				return
			}
		})
	}

	if exitCode > 0 && email {

		wg.Go(func() {
			subject := "Command failed on " + hostName
			body := fmt.Sprintf(
				"%12s%q\n%12s%s\n%12s%s",
				"Command: ", command,
				"Start time: ", startTime.Format(MAILDATE),
				"Stop time: ", stopTime.Format(MAILDATE),
			)
			err := SendLogEmail(subject, body, stdOutFile, stdErrFile)
			if err != nil {
				fmt.Println(err)
			}
		})
	}

	wg.Wait()

	return nil
}
