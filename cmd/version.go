/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  "Print the version number of errwrapper",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Errwrapper v" + Version)
	},
}