/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"errors"
	"io"
	"os"
	"sync"
)

type Exit struct{ Code int }

func HandleExit() {
	if e := recover(); e != nil {
		if exit, ok := e.(Exit); ok == true {
			os.Exit(exit.Code)
		}
		panic(e)
	}
}

func GetEnvVar(variable string) (string, error) {
	v := os.Getenv(variable)
	if v == "" {
		return "", errors.New("Variable " + variable + " is empty. Exiting.")
	}

	return v, nil
}

func Tee(in io.Reader, wg *sync.WaitGroup, out ...string) error {
	defer wg.Done()

	var fileDescriptors []io.Writer

	for _, element := range out {
		fileDescriptor, err := os.OpenFile(element, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return errors.New("Failed to open file.")
		}
		defer func() error {
			err := fileDescriptor.Close()
			if err != nil {
				return errors.New("Failed to close file.")
			}

			return nil
		}() // close each fd after we finish reading
		fileDescriptors = append(fileDescriptors, fileDescriptor)
	}

	if !Quiet {
		fileDescriptors = append(fileDescriptors, os.Stdout)
	}

	writer := io.MultiWriter(fileDescriptors...)

	buf := make([]byte, 256)

	_, err := io.CopyBuffer(writer, in, buf)
	if err != nil {
		return errors.New("Failed to write logs.")
	}

	return nil
}
