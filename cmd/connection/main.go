package connect

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/orginux/sql-cd/cmd/logging"
)

// Clickhouse connects to the ClickHouse database via native port
func Clickhouse(destHost string, destPort int, username, password string, dbVerboseMode bool) (context.Context, clickhouse.Conn, error) {

	var (
		ctx      context.Context
		conn     clickhouse.Conn
		attemppt uint

		clickhouseURL = fmt.Sprintf("%s:%d", destHost, destPort)
	)

	for {
		var err error
		conn, err = clickhouse.Open(&clickhouse.Options{
			Addr: []string{clickhouseURL},
			Auth: clickhouse.Auth{
				Database: "default",
				Username: username,
				Password: password,
			},
			Debug:           dbVerboseMode,
			DialTimeout:     time.Second,
			MaxOpenConns:    10,
			MaxIdleConns:    5,
			ConnMaxLifetime: time.Hour,
			Compression:     &clickhouse.Compression{Method: clickhouse.CompressionLZ4},
		})
		if err != nil {
			return nil, nil, err
		}

		ctx = clickhouse.Context(context.Background(), clickhouse.WithSettings(clickhouse.Settings{
			"max_block_size": 10,
		}), clickhouse.WithProgress(func(p *clickhouse.Progress) {
			fmt.Println("progress: ", p)
		}), clickhouse.WithProfileInfo(func(p *clickhouse.ProfileInfo) {
			fmt.Println("profile info: ", p)
		}))

		if err := conn.Ping(ctx); err != nil {
			if exception, ok := err.(*clickhouse.Exception); ok {
				fmt.Printf("Catch exception [%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
			}
			attemppt++
			logging.Error.Printf("Connection attempt #%d failed with error: %v", attemppt, err)
			time.Sleep(5 * time.Second)
			continue
		}

		logging.Info.Printf("Connection established: %s\n", clickhouseURL)
		break
	}
	return ctx, conn, nil
}
