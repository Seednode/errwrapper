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
	defer HandleExit()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, os.Interrupt)

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Home directory not found.")
		panic(Exit{1})
	}

	envFile := homeDirectory + "/.config/errwrapper/.env"
	err = godotenv.Load(envFile)
	if err != nil {
		fmt.Printf("Failed to load env file %q.", envFile)
		panic(Exit{1})
	}

	hostName, err := os.Hostname()
	if err != nil {
		fmt.Println("Hostname not found.")
		panic(Exit{1})
	}

	startTime := time.Now()
	stdOutFile, stdErrFile, exitCode, err := LogCommand(arguments)
	if err != nil {
		fmt.Println(err)
	}
	stopTime := time.Now()

	command := strings.Join(arguments, " ")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		dbType, err := GetEnvVar("ERRWRAPPER_DB_TYPE", false)
		if err != nil {
			fmt.Println(err)
			return
		}

		if dbType != "cockroachdb" && dbType != "postgresql" {
			fmt.Println("invalid database type specified")
			return
		}

		databaseURL, err := GetDatabaseURL(dbType)
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
			err := SendLogEmail(subject, body, stdOutFile, stdErrFile)
			if err != nil {
				fmt.Println(err)
			}
		}()
	}

	wg.Wait()
}
