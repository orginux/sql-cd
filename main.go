package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/orginux/sql-cd/cmd/apply"
	connect "github.com/orginux/sql-cd/cmd/connection"
	git "github.com/orginux/sql-cd/cmd/git"
	logging "github.com/orginux/sql-cd/cmd/logging"
)

var (
	destHost                            string
	destPort                            int
	username, password                  string
	gitURL, gitBranch, gitPath, gitDest string
	workDir                             string
	runAsDaemon                         bool
	timeout                             int
)

func init() {
	// connection
	flag.StringVar(&destHost, "host", "localhost", "The ClickHouse server name. You can use either the name or the IPv4 or IPv6 address")
	flag.IntVar(&destPort, "port", 9000, "The native ClickHouse port to connect to")
	flag.StringVar(&username, "username", "default", "The username. Default value: default")
	flag.StringVar(&password, "password", "", "The password. Default value: empty string")

	// git
	flag.StringVar(&gitURL, "git-url", "", "URL of git repo with SQL queries")
	flag.StringVar(&gitBranch, "git-branch", "main", "Branch of git repo to use for SQL queries")
	flag.StringVar(&gitPath, "git-path", "", "Path within git repo to locate SQL queries")
	flag.StringVar(&workDir, "work-dir", "/tmp/ufo/", "Local path for repo with SQL queries")
	// flag.StringVar(&gitDest, "git-dest", "", "local path for repo with SQL queries")

	// daemon
	flag.BoolVar(&runAsDaemon, "daemon", false, "Run as daemon")
	flag.IntVar(&timeout, "timeout", 60, "Global command timeout")

	flag.Parse()
}

func checkErr(err error, runAsDaemon bool) {
	if err != nil {
		// Daemon doesn't exit when catch an error
		if !runAsDaemon {
			log.Fatal(err)
		}
		logging.Error.Println(err)
	}

}

func main() {
	// Set work dir
	subPath := strings.ReplaceAll(gitURL, "https://", "")
	gitDest := filepath.Join(workDir, subPath, gitBranch)
	queriesDir := filepath.Join(gitDest, gitPath)

	for {
		// Connect to ClickHouse
		ctx, conn, err := connect.Clickhouse(destHost, destPort, username, password)
		checkErr(err, runAsDaemon)

		// Clone project
		err = git.Clone(gitDest, gitURL, gitBranch)
		checkErr(err, runAsDaemon)

		// Apply SQL files
		apply.QueriesFromDir(ctx, conn, queriesDir)

		// Close connection
		conn.Close()
		logging.Info.Println("Connection closed")

		// Exit if run once
		if !runAsDaemon {
			os.Exit(0)
		}

		// Wait before next iteration
		time.Sleep(time.Duration(timeout) * time.Second)
	}
}
