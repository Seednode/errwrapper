/*
Copyright © 2025 Seednode <seednode@seedno.de>
*/

package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ReleaseVersion string = "1.0.0"
)

var (
	databaseType     string
	databaseHost     string
	databasePort     string
	databaseUser     string
	databasePass     string
	databaseName     string
	databaseTable    string
	databaseSslMode  string
	databaseRootCert string
	databaseSslCert  string
	databaseSslKey   string
	loggingDirectory string
	mailServer       string
	mailPort         string
	mailFrom         string
	mailTo           string
	mailUser         string
	mailPass         string
	database         bool
	email            bool
	stdOut           bool
	verbose          bool
)

func main() {
	cmd := &cobra.Command{
		Use:   "errwrapper <command>",
		Short: "Runs a command, logging output to a file and a database, emailing if the command fails.",
		Args:  cobra.MinimumNArgs(1),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			initializeConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunCommand(args)
		},
	}

	cmd.Flags().StringVar(&databaseType, "db-type", "", "database type to connect to")
	cmd.Flags().StringVar(&databaseHost, "db-host", "", "database host to connect to")
	cmd.Flags().StringVar(&databasePort, "db-port", "", "database port to connect to")
	cmd.Flags().StringVar(&databaseUser, "db-user", "", "database user to connect as")
	cmd.Flags().StringVar(&databasePass, "db-pass", "", "database password to connect with")
	cmd.Flags().StringVar(&databaseName, "db-name", "", "database name to connect to")
	cmd.Flags().StringVar(&databaseTable, "db-table", "", "database table to query")
	cmd.Flags().StringVar(&databaseSslMode, "db-ssl-mode", "", "database ssl connection mode")
	cmd.Flags().StringVar(&databaseRootCert, "db-root-cert", "", "database ssl root certificate path")
	cmd.Flags().StringVar(&databaseSslCert, "db-ssl-cert", "", "database ssl connection certificate path")
	cmd.Flags().StringVar(&databaseSslKey, "db-ssl-key", "", "database ssl connection key path")
	cmd.Flags().StringVarP(&loggingDirectory, "logging-directory", "l", "", "directory to log to (defaults to $HOME/errwrapper)")
	cmd.Flags().StringVar(&mailServer, "mail-server", "", "mailserver to use for error notifications")
	cmd.Flags().StringVar(&mailPort, "mail-port", "", "smtp port for mailserver")
	cmd.Flags().StringVar(&mailFrom, "mail-from", "", "from address to use for error notifications")
	cmd.Flags().StringVar(&mailTo, "mail-to", "", "recipient for error notifications")
	cmd.Flags().StringVar(&mailUser, "mail-user", "", "username for smtp account")
	cmd.Flags().StringVar(&mailPass, "mail-pass", "", "password for smtp account")
	cmd.Flags().BoolVarP(&database, "database", "d", false, "log command info to database")
	cmd.Flags().BoolVarP(&email, "email", "e", false, "send email on error")
	cmd.Flags().BoolVarP(&stdOut, "stdout", "s", false, "log output to stdout as well as a file")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "display environment variables on start")

	cmd.Flags().SetInterspersed(true)

	cmd.CompletionOptions.HiddenDefaultCmd = true

	cmd.SilenceErrors = true
	cmd.SetHelpCommand(&cobra.Command{
		Hidden: true,
	})

	cmd.SetVersionTemplate("errwrapper v{{.Version}}\n")
	cmd.Version = ReleaseVersion

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetEnvPrefix("errwrapper")

	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	v.AutomaticEnv()

	bindFlags(cmd, v)

	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		configName := strings.ReplaceAll(f.Name, "-", "_")

		if !f.Changed && v.IsSet(configName) {
			val := v.Get(configName)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
