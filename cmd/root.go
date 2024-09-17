/*
Copyright © 2024 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	ReleaseVersion string = "0.4.0"
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

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "errwrapper <command>",
		Short: "Runs a command, logging output to a file and a database, emailing if the command fails.",
		Args:  cobra.MinimumNArgs(1),
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := RunCommand(args)
			if err != nil {
				return err
			}

			return nil
		},
	}

	rootCmd.Flags().StringVar(&DatabaseType, "db-type", "", "database type to connect to")
	rootCmd.Flags().StringVar(&DatabaseHost, "db-host", "", "database host to connect to")
	rootCmd.Flags().StringVar(&DatabasePort, "db-port", "", "database port to connect to")
	rootCmd.Flags().StringVar(&DatabaseUser, "db-user", "", "database user to connect as")
	rootCmd.Flags().StringVar(&DatabasePass, "db-pass", "", "database password to connect with")
	rootCmd.Flags().StringVar(&DatabaseName, "db-name", "", "database name to connect to")
	rootCmd.Flags().StringVar(&DatabaseTable, "db-table", "", "database table to query")
	rootCmd.Flags().StringVar(&DatabaseSslMode, "db-ssl-mode", "", "database ssl connection mode")
	rootCmd.Flags().StringVar(&DatabaseRootCert, "db-root-cert", "", "database ssl root certificate path")
	rootCmd.Flags().StringVar(&DatabaseSslCert, "db-ssl-cert", "", "database ssl connection certificate path")
	rootCmd.Flags().StringVar(&DatabaseSslKey, "db-ssl-key", "", "database ssl connection key path")
	rootCmd.Flags().StringVarP(&LoggingDirectory, "logging-directory", "l", "", "directory to log to (defaults to $HOME/errwrapper)")
	rootCmd.Flags().StringVar(&MailServer, "mail-server", "", "mailserver to use for error notifications")
	rootCmd.Flags().StringVar(&MailPort, "mail-port", "", "smtp port for mailserver")
	rootCmd.Flags().StringVar(&MailFrom, "mail-from", "", "from address to use for error notifications")
	rootCmd.Flags().StringVar(&MailTo, "mail-to", "", "recipient for error notifications")
	rootCmd.Flags().StringVar(&MailUser, "mail-user", "", "username for smtp account")
	rootCmd.Flags().StringVar(&MailPass, "mail-pass", "", "password for smtp account")
	rootCmd.Flags().StringVar(&TimeZone, "timezone", "", "timezone to use")
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

	return rootCmd
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()

	v.SetConfigName("config")

	v.SetConfigType("yaml")

	v.AddConfigPath("/etc/errwrapper/")
	v.AddConfigPath("$HOME/.config/errwrapper")
	v.AddConfigPath(".")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

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
