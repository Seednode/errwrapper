/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

const LOGDATE string = "2006-01-02"

func CreateLoggingDirectory() string {
	now := time.Now()
	currentDate := now.Format(LOGDATE)

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	loggingDirectory := homeDirectory + "/logs/" + currentDate

	err = os.MkdirAll(loggingDirectory, 0755)
	if err != nil {
		panic(err)
	}

	return loggingDirectory
}

func LogCommand(pidFile *os.File, arguments []string) (string, string, int) {
	timeStamp := fmt.Sprint(time.Now().UnixMicro())
	loggingDirectory := CreateLoggingDirectory()
	loggingPrefix := loggingDirectory + "/" + timeStamp + "_" + filepath.Base(arguments[0])
	stdOutFile := loggingPrefix + "_out.log"
	stdErrFile := loggingPrefix + "_err.log"

	cmd := exec.Command(arguments[0], arguments[1:]...)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		panic(err)
	}

	err = cmd.Start()
	if err != nil {
		panic(err)
	}

	pid := strconv.Itoa(cmd.Process.Pid)
	_, err = pidFile.Write([]byte(pid))
	if err != nil {
		panic(err)
	}

	err = pidFile.Close()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go Tee(stdOut, &wg, stdOutFile)

	wg.Add(1)
	go Tee(stdErr, &wg, stdErrFile)

	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}

	exitCode := cmd.ProcessState.ExitCode()
	if err != nil {
		panic(err)
	}

	RemovePIDFile()

	return stdOutFile, stdErrFile, exitCode
}
