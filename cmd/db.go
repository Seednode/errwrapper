/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/jackc/pgx/v5"
)

const DBDATE string = "2006-01-02T15:04:05.000000000-07:00"

func GetDatabaseURL(dbType string) (string, error) {
	var url strings.Builder

	host, err := GetEnvVar("ERRWRAPPER_DB_HOST", DatabaseHost, false)
	if err != nil {
		return "", err
	}
	url.WriteString("host=" + host)

	port, err := GetEnvVar("ERRWRAPPER_DB_PORT", DatabasePort, false)
	if err != nil {
		return "", err
	}
	url.WriteString(" port=" + port)

	user, err := GetEnvVar("ERRWRAPPER_DB_USER", DatabaseUser, false)
	if err != nil {
		return "", err
	}
	url.WriteString(" user=" + user)

	if dbType == "postgresql" {
		pass, err := GetEnvVar("ERRWRAPPER_DB_PASS", DatabasePass, true)
		if err != nil {
			return "", err
		}
		url.WriteString(" password=" + pass)
	}

	database, err := GetEnvVar("ERRWRAPPER_DB_NAME", DatabaseName, false)
	if err != nil {
		return "", err
	}
	url.WriteString(" dbname=" + database)

	sslMode, err := GetEnvVar("ERRWRAPPER_DB_SSL_MODE", DatabaseSslMode, false)
	if err != nil {
		return "", err
	}
	url.WriteString(" sslmode=" + sslMode)

	if dbType == "cockroachdb" {
		sslRootCert, err := GetEnvVar("ERRWRAPPER_DB_ROOT_CERT", DatabaseRootCert, false)
		if err != nil {
			return "", err
		}
		url.WriteString(" sslrootcert=" + sslRootCert)

		sslClientKey, err := GetEnvVar("ERRWRAPPER_DB_SSL_KEY", DatabaseSslKey, false)
		if err != nil {
			return "", err
		}
		url.WriteString(" sslkey=" + sslClientKey)

		sslClientCert, err := GetEnvVar("ERRWRAPPER_DB_SSL_CERT", DatabaseSslCert, false)
		if err != nil {
			return "", err
		}
		url.WriteString(" sslcert=" + sslClientCert)
	}

	return url.String(), nil
}

func CreateSQLStatement(startTime, stopTime time.Time, hostName, command string, exitCode int) (string, error) {
	tableName, err := GetEnvVar("ERRWRAPPER_DB_TABLE", DatabaseTable, false)
	if err != nil {
		return "", err
	}

	fields := "startTime, stopTime, hostName, commandName, exitCode"
	values := [5]string{
		startTime.Format(DBDATE),
		stopTime.Format(DBDATE),
		hostName,
		command,
		strconv.Itoa(exitCode),
	}

	var data string
	for value := 0; value < len(values); value++ {
		data += fmt.Sprintf("'%s', ", values[value])
	}

	dataToInsert := strings.TrimSuffix(data, ", ")

	statement := "INSERT INTO " + tableName + "(" + fields + ") VALUES (" + dataToInsert + ");"

	return statement, nil
}

func WriteToDatabase(databaseURL, sqlStatement string) error {
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		return errors.New("failed to connect to database")
	}

	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}(conn, context.Background())

	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return InsertRows(context.Background(), tx, sqlStatement)
	})
	if err != nil {
		return err
	}

	return nil
}

func InsertRows(ctx context.Context, tx pgx.Tx, statement string) error {
	_, err := tx.Exec(ctx, statement)

	if err != nil {
		return errors.New("failed to execute database transaction")
	}

	return nil
}
