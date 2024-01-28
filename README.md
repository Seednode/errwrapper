## About
A tool to log the execution of commands to a file, email, a database, or any combination thereof.

Feature requests, code criticism, bug reports, general chit-chat, and unrelated angst accepted at `errwrapper@seedno.de`.

Static binary builds available [here](https://cdn.seedno.de/builds/errwrapper).

x86_64 and ARM Docker images of latest version: `oci.seedno.de/seednode/errwrapper:latest`.

Dockerfile available [here](https://git.seedno.de/seednode/errwrapper/raw/branch/master/docker/Dockerfile).

## Creating the table
The command database is designed to be viewed via the accompanying [commands](https://git.seedno.de/seednode/commands) tool, which should connect to the same database.

In this example, I'll be using the wonderful [usql](https://github.com/xo/usql) client.

### Connect to the database:
To connect, run the following, replacing the variables with their corresponding values:

`usql postgres://${ERRWRAPPER_DB_USER}@${ERRWRAPPER_DB_HOST}:${ERRWRAPPER_DB_PORT}/${ERRWRAPPER_DB_NAME}`

You should then be at a SQL prompt that looks something like the following:

`pg:errwrapper@errwrapper-db/logging=>`

### Create logging table
To create a table with the proper structure, run the following (as always, adjusting variables as needed):
```
CREATE TABLE ${ERRWRAPPER_DB_TABLE} (
	id SERIAL PRIMARY KEY,
	starttime timestamp NOT NULL,
	stoptime timestamp NOT NULL,
	hostname varchar NOT NULL,
	commandname varchar NOT NULL,
	exitcode int NOT NULL
);
```

## Configure the tool
The following environment variables are used to configure errwrapper (all values provided are just examples):
```
ERRWRAPPER_DB_TYPE=postgresql
ERRWRAPPER_DB_HOST=errwrapper-db
ERRWRAPPER_DB_PORT=5432
ERRWRAPPER_DB_USER=errwrapper
ERRWRAPPER_DB_PASS=changeme
ERRWRAPPER_DB_NAME=logging
ERRWRAPPER_DB_TABLE=logging
ERRWRAPPER_DB_SSL_MODE=disable
ERRWRAPPER_MAIL_SERVER=smtp.fake.example
ERRWRAPPER_MAIL_PORT=465
ERRWRAPPER_MAIL_FROM=me@fake.example
ERRWRAPPER_MAIL_TO=you@fake.example
ERRWRAPPER_MAIL_USER=errwrapper@fake.example
ERRWRAPPER_MAIL_PASS=changemetoo
TZ=America/Chicago
```

## Usage output
Alternatively, you can configure errwrapper using command-line flags.
```
Runs a command, logging output to a file and a database, emailing if the command fails.

Usage:
  errwrapper <command> [flags]

Flags:
  -d, --database                    log command info to database
      --database-host string        database host to connect to
      --database-name string        database name to connect to
      --database-pass string        database password to connect with
      --database-port string        database port to connect to
      --database-root-cert string   database ssl root certificate path
      --database-ssl-cert string    database ssl connection certificate path
      --database-ssl-key string     database ssl connection key path
      --database-ssl-mode string    database ssl connection mode
      --database-table string       database table to query
      --database-type string        database type to connect to
      --database-user string        database user to connect as
  -e, --email                       send email on error
  -h, --help                        help for errwrapper
  -l, --logging-directory string    directory to log to (defaults to $HOME/errwrapper)
      --mail-from string            from address to use for error notifications
      --mail-pass string            password for smtp account
      --mail-port string            smtp port for mailserver
      --mail-server string          mailserver to use for error notifications
      --mail-to string              recipient for error notifications
      --mail-user string            username for smtp account
  -s, --stdout                      log output to stdout as well as a file
      --time-zone string            timezone to use
  -v, --verbose                     display environment variables on start
```