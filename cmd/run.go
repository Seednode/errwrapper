/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func RunCommand(arguments []string) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, os.Interrupt)

	pidFile := CreatePIDFile()
	defer func() {
		err := pidFile.Close()
		if err != nil {
			panic(err)
		}

		RemovePIDFile()
	}()

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	envFile := homeDirectory + "/.config/errwrapper/.env"
	err = godotenv.Load(envFile)
	if err != nil {
		panic(err)
	}

	startTime := time.Now()
	stdOutFile, stdErrFile, exitCode := LogCommand(pidFile, arguments)
	stopTime := time.Now()

	hostName, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	command := strings.Join(arguments, " ")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		databaseURL := GetDatabaseURL()
		sqlStatement := CreateSQLStatement(startTime, stopTime, hostName, command, exitCode)
		WriteToDatabase(databaseURL, sqlStatement)
	}()

	if exitCode > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			subject := "Command failed on " + hostName
			body := fmt.Sprintf(
				"%12s%q\n%12s%s\n%12s%s",
				"Command: ", command,
				"Start time: ", startTime.Format(MAILDATE),
				"Stop time: ", stopTime.Format(MAILDATE),
			)
			SendLogEmail(subject, body, stdOutFile, stdErrFile)
		}()
	}

	wg.Wait()
}
