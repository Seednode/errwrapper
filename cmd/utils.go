/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"io"
	"os"
	"sync"
)

func GetEnvVar(variable string) string {
	v := os.Getenv(variable)
	if v == "" {
		fmt.Println("Variable " + variable + " is empty. Exiting.")
		os.Exit(1)
	}

	return v
}

func Tee(in io.Reader, wg *sync.WaitGroup, out ...string) {
	defer wg.Done()

	var fileDescriptors []io.Writer

	for _, element := range out {
		fileDescriptor, err := os.OpenFile(element, os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			panic(err)
		}
		defer func() {
			err := fileDescriptor.Close()
			if err != nil {
				panic(err)
			}
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
		panic(err)
	}
}
