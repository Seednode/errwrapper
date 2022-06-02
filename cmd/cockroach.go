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
	"sync"
	"time"

	"github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgx"
	pgx "github.com/jackc/pgx/v4"
)

const DBDATE string = "2006-01-02T15:04:05.000000000-07:00"

func GetDatabaseURL() (string, error) {
	hostVar, err := GetEnvVar("ERRWRAPPER_DB_HOST")
	if err != nil {
		return "", err
	}
	host := "host=" + hostVar

	portVar, err := GetEnvVar("ERRWRAPPER_DB_PORT")
	if err != nil {
		return "", err
	}
	port := " port=" + portVar

	userVar, err := GetEnvVar("ERRWRAPPER_DB_USER")
	if err != nil {
		return "", err
	}
	user := " user=" + userVar

	databaseVar, err := GetEnvVar("ERRWRAPPER_DB_NAME")
	if err != nil {
		return "", err
	}
	database := " dbname=" + databaseVar

	sslModeVar, err := GetEnvVar("ERRWRAPPER_DB_SSL_MODE")
	if err != nil {
		return "", err
	}
	sslMode := " sslmode=" + sslModeVar

	sslRootCertVar, err := GetEnvVar("ERRWRAPPER_DB_ROOT_CERT")
	if err != nil {
		return "", err
	}
	sslRootCert := " sslrootcert=" + sslRootCertVar

	sslClientKeyVar, err := GetEnvVar("ERRWRAPPER_DB_SSL_KEY")
	if err != nil {
		return "", err
	}
	sslClientKey := " sslkey=" + sslClientKeyVar

	sslClientCertVar, err := GetEnvVar("ERRWRAPPER_DB_SSL_CERT")
	if err != nil {
		return "", err
	}
	sslClientCert := " sslcert=" + sslClientCertVar

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

	return connection, nil
}

func CreateSQLStatement(startTime, stopTime time.Time, hostName, command string, exitCode int) (string, error) {
	tableName, err := GetEnvVar("ERRWRAPPER_DB_TABLE")
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
		return errors.New("Failed to connect to database.")
	}

	defer conn.Close(context.Background())

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
		return errors.New("Failed to execute database transaction.")
	}

	return nil
}

func LogToDatabase(startTime, stopTime time.Time, hostName, command string, exitCode int, wg *sync.WaitGroup) {
	defer wg.Done()

	databaseURL, err := GetDatabaseURL()
	if err != nil {
		fmt.Println(err)
		panic(Exit{1})
	}

	sqlStatement, err := CreateSQLStatement(startTime, stopTime, hostName, command, exitCode)
	if err != nil {
		fmt.Println(err)
		panic(Exit{1})
	}

	err = WriteToDatabase(databaseURL, sqlStatement)
	if err != nil {
		fmt.Println(err)
		panic(Exit{1})
	}
}
