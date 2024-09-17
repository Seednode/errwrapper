/*
Copyright © 2024 Seednode <seednode@seedno.de>
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

func GetDatabaseURL() (string, error) {
	var url strings.Builder

	url.WriteString("host=" + DatabaseHost)
	url.WriteString(" port=" + DatabasePort)
	url.WriteString(" user=" + DatabaseUser)

	if DatabaseType == "postgresql" {
		url.WriteString(" password=" + DatabasePass)
	}

	url.WriteString(" dbname=" + DatabaseName)
	url.WriteString(" sslmode=" + DatabaseSslMode)

	if DatabaseType == "cockroachdb" {
		url.WriteString(" sslrootcert=" + DatabaseRootCert)
		url.WriteString(" sslkey=" + DatabaseSslKey)
		url.WriteString(" sslcert=" + DatabaseSslCert)
	}

	return url.String(), nil
}

func CreateSQLStatement(startTime, stopTime time.Time, hostName, command string, exitCode int) (string, error) {
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

	statement := "INSERT INTO " + DatabaseTable + "(" + fields + ") VALUES (" + dataToInsert + ");"

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
