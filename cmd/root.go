/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	ReleaseVersion string = "0.3.2"
)

var (
	DatabaseType     string
	DatabaseHost     string
	DatabasePort     string
	DatabaseUser     string
	DatabasePass     string
	DatabaseName     string
	DatabaseTable    string
	DatabaseSslMode  string
	DatabaseRootCert string
	DatabaseSslCert  string
	DatabaseSslKey   string
	LoggingDirectory string
	MailServer       string
	MailPort         string
	MailFrom         string
	MailTo           string
	MailUser         string
	MailPass         string
	TimeZone         string
	Database         bool
	Email            bool
	StdOut           bool
	Verbose          bool
)

var rootCmd = &cobra.Command{
	Use:   "errwrapper <command>",
	Short: "Runs a command, logging output to a file and a database, emailing if the command fails.",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := RunCommand(args)
		if err != nil {
			return err
		}

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVar(&DatabaseType, "database-type", "", "database type to connect to")
	rootCmd.Flags().StringVar(&DatabaseHost, "database-host", "", "database host to connect to")
	rootCmd.Flags().StringVar(&DatabasePort, "database-port", "", "database port to connect to")
	rootCmd.Flags().StringVar(&DatabaseUser, "database-user", "", "database user to connect as")
	rootCmd.Flags().StringVar(&DatabasePass, "database-pass", "", "database password to connect with")
	rootCmd.Flags().StringVar(&DatabaseName, "database-name", "", "database name to connect to")
	rootCmd.Flags().StringVar(&DatabaseTable, "database-table", "", "database table to query")
	rootCmd.Flags().StringVar(&DatabaseSslMode, "database-ssl-mode", "", "database ssl connection mode")
	rootCmd.Flags().StringVar(&DatabaseRootCert, "database-root-cert", "", "database ssl root certificate path")
	rootCmd.Flags().StringVar(&DatabaseSslCert, "database-ssl-cert", "", "database ssl connection certificate path")
	rootCmd.Flags().StringVar(&DatabaseSslKey, "database-ssl-key", "", "database ssl connection key path")
	rootCmd.Flags().StringVarP(&LoggingDirectory, "logging-directory", "l", "", "directory to log to (defaults to $HOME/errwrapper)")
	rootCmd.Flags().StringVar(&MailServer, "mail-server", "", "mailserver to use for error notifications")
	rootCmd.Flags().StringVar(&MailPort, "mail-port", "", "smtp port for mailserver")
	rootCmd.Flags().StringVar(&MailFrom, "mail-from", "", "from address to use for error notifications")
	rootCmd.Flags().StringVar(&MailTo, "mail-to", "", "recipient for error notifications")
	rootCmd.Flags().StringVar(&MailUser, "mail-user", "", "username for smtp account")
	rootCmd.Flags().StringVar(&MailPass, "mail-pass", "", "password for smtp account")
	rootCmd.Flags().StringVar(&TimeZone, "time-zone", "", "timezone to use")
	rootCmd.Flags().BoolVarP(&Database, "database", "d", false, "log command info to database")
	rootCmd.Flags().BoolVarP(&Email, "email", "e", false, "send email on error")
	rootCmd.Flags().BoolVarP(&StdOut, "stdout", "s", false, "log output to stdout as well as a file")
	rootCmd.Flags().BoolVarP(&Verbose, "verbose", "v", false, "display environment variables on start")

	rootCmd.Flags().SetInterspersed(true)

	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	rootCmd.SilenceErrors = true
	rootCmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	rootCmd.SetVersionTemplate("errwrapper v{{.Version}}\n")
	rootCmd.Version = ReleaseVersion
}
