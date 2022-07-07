package apply

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/orginux/sql-cd/cmd/logging"
)

// QueriesFromDir applies all queries from a directory
func QueriesFromDir(ctx context.Context, conn clickhouse.Conn, queriesDir string) error {

	logging.Debug.Printf("queries directory: %s\n", queriesDir)
	queryFiles, err := ioutil.ReadDir(queriesDir)
	if err != nil {
		logging.Error.Println(err)
		return err
	}

	var wg sync.WaitGroup

	logging.Debug.Printf("queryFiles: %v\n", queryFiles)
	for _, fileFromDir := range queryFiles {
		wg.Add(1)

		queryFile := fileFromDir.Name()
		logging.Debug.Printf("File Name: %s", queryFile)
		go func() {
			defer wg.Done()
			err := applyFile(ctx, conn, filepath.Join(queriesDir, queryFile))
			if err != nil {
				// Ignore query errors
				logging.Error.Printf("%s: %v\n", queryFile, err)
			}
		}()
	}
	// Wait for all requests from the directory to complete
	wg.Wait()

	// Remove queries after applaing
	logging.Debug.Println("Remove directory")
	os.RemoveAll(queriesDir)

	return nil
}

// applyFile applies all queries from a file
func applyFile(ctx context.Context, conn clickhouse.Conn, queryFile string) error {
	content, err := ioutil.ReadFile(queryFile)
	if err != nil {
		return err
	}

	queries := strings.Split(string(content), ";")

	// Print count of queries for each file
	logging.Info.Printf("File: %s, queries: %d", queryFile, len(queries))

	for _, query := range queries {
		if len(query) > 0 && query != "\n" {
			err = conn.Exec(ctx, query)
			if err != nil {
				return fmt.Errorf("Query: %s, Error: %s", query, err)
			}
		}
	}
	return nil
}
