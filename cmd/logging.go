/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const LOGDATE string = "2006-01-02"

func CreateLoggingDirectory() (string, error) {
	now := time.Now()
	currentDate := now.Format(LOGDATE)

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("home directory not found")
	}

	loggingDirectory := homeDirectory + "/logs/" + currentDate

	err = os.MkdirAll(loggingDirectory, 0755)
	if err != nil {
		return "", errors.New("failed to create logging directory")
	}

	return loggingDirectory, nil
}

func LogCommand(pidFile *os.File, arguments []string) (string, string, int, error) {
	timeStamp := fmt.Sprint(time.Now().UnixMicro())
	loggingDirectory, err := CreateLoggingDirectory()
	if err != nil {
		return "", "", 0, err
	}
	loggingPrefix := loggingDirectory + "/" + timeStamp + "_" + filepath.Base(arguments[0])
	stdOutFile := loggingPrefix + "_out.log"
	stdErrFile := loggingPrefix + "_err.log"

	cmd := exec.Command(arguments[0], arguments[1:]...)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return "", "", 0, errors.New("failed to allocate pipe for stdout")
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		return "", "", 0, errors.New("failed to allocate pipe for stderr")
	}

	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
	}

	pid := strconv.Itoa(cmd.Process.Pid)
	_, err = pidFile.Write([]byte(pid))
	if err != nil {
		return "", "", 0, errors.New("failed to write pid file")
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		err := Tee(stdOut, &wg, stdOutFile)
		if err != nil {
			fmt.Println(err)
		}
	}()

	wg.Add(1)
	go func() {
		err := Tee(stdErr, &wg, stdErrFile)
		if err != nil {
			fmt.Println(err)
		}
	}()

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		return "", "", 0, err
	}

	exitCode := cmd.ProcessState.ExitCode()

	return stdOutFile, stdErrFile, exitCode, nil
}
