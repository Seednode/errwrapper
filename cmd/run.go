/*
Copyright © 2023 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func RunCommand(arguments []string) error {
	defer HandleExit()

	timezone, err := GetEnvVar("TZ", TimeZone, false)
	if err != nil {
		timezone = "UTC"
	}

	time.Local, err = time.LoadLocation(timezone)
	if err != nil {
		return err
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, os.Interrupt)

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return errors.New("home directory not found")
	}

	envFile := filepath.Join(homeDirectory, ".config", "errwrapper", ".env")
	err = godotenv.Load(envFile)
	if err != nil {
		return fmt.Errorf("failed to load env file %q", envFile)
	}

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

	if Database {
		wg.Add(1)

		go func() {
			defer wg.Done()

			dbType, err := GetEnvVar("ERRWRAPPER_DB_TYPE", DatabaseType, false)
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
	}

	if exitCode > 0 && Email {
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

	return nil
}
