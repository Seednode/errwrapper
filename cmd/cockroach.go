/*
Copyright © 2022 Seednode <seednode@seedno.de>
*/

package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	pgx "github.com/jackc/pgx/v4"
)

const DBDATE string = "2006-01-02T15:04:05.000000000-07:00"

func GetDatabaseURL() string {
	host := "host=" + GetEnvVar("ERRWRAPPER_DB_HOST")
	port := " port=" + GetEnvVar("ERRWRAPPER_DB_PORT")
	user := " user=" + GetEnvVar("ERRWRAPPER_DB_USER")
	database := " dbname=" + GetEnvVar("ERRWRAPPER_DB_NAME")
	sslMode := " sslmode=" + GetEnvVar("ERRWRAPPER_DB_SSL_MODE")
	sslRootCert := " sslrootcert=" + GetEnvVar("ERRWRAPPER_DB_ROOT_CERT")
	sslClientKey := " sslkey=" + GetEnvVar("ERRWRAPPER_DB_SSL_KEY")
	sslClientCert := " sslcert=" + GetEnvVar("ERRWRAPPER_DB_SSL_CERT")

	connection := fmt.Sprint(
		host,
		port,
		user,
		database,
		sslMode,
		sslRootCert,
		sslClientKey,
		sslClientCert,
	)

	return connection
}

func CreateSQLStatement(startTime, stopTime time.Time, hostName, command string, exitCode int) string {
	tableName := GetEnvVar("ERRWRAPPER_DB_TABLE")
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

	return statement
}

func WriteToDatabase(databaseURL, sqlStatement string) {
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		panic(err)
	}

	defer conn.Close(context.Background())

	err = crdbpgx.ExecuteTx(context.Background(), conn, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return InsertRows(context.Background(), tx, sqlStatement)
	})
	if err != nil {
		panic(err)
	}
}

func InsertRows(ctx context.Context, tx pgx.Tx, statement string) error {
	_, err := tx.Exec(ctx, statement)

	return err
}
