/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func CheckPIDFile() (string, bool) {
	homeDirectory, _ := os.UserHomeDir()
	pidFile := homeDirectory + "/pids/" + filepath.Base(os.Args[1]) + ".pid"
	if _, err := os.Stat(pidFile); !errors.Is(err, os.ErrNotExist) {
		return pidFile, true
	}
	return pidFile, false
}

func CreatePIDFile() *os.File {
	pidFile, exists := CheckPIDFile()
	if exists {
		fmt.Println("Pidfile exists. Not running.")
		os.Exit(1)
	}
	filePtr, err := os.Create(pidFile)
	if err != nil {
		panic(err)
	}
	return filePtr
}

func RemovePIDFile() {
	pidfile, exists := CheckPIDFile()
	if exists {
		err := os.Remove(pidfile)
		if err != nil {
			panic(err)
		}
	}
}
