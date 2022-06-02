/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version string = "0.1"
var Quiet bool

var rootCmd = &cobra.Command{
	Use:   "errwrapper <command>",
	Short: "Runs a command, logging output to a file and a cockroachdb table, emailing if the command fails.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		RunCommand(args)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		panic(Exit{1})
	}
}

func init() {
	rootCmd.Flags().BoolVarP(&Quiet, "quiet", "q", false, "only write output to file")
	rootCmd.Flags().SetInterspersed(false)
}
