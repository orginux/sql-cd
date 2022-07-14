package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/orginux/sql-cd/cmd/apply"
	"github.com/orginux/sql-cd/cmd/config"
	connect "github.com/orginux/sql-cd/cmd/connection"
	git "github.com/orginux/sql-cd/cmd/git"
	logging "github.com/orginux/sql-cd/cmd/logging"
)

// DataBase variables
var (
	dbHost                 string
	dbPort                 int
	dbUsername, dbPassword string
	dbVerboseMode          bool
)

// Git variables
var (
	gitURL, gitBranch, gitPrivateKeyFile string
	gitPath, gitDest                     string
	workDir                              string
)

// Service variables
var (
	runAsDaemon bool
	timeout     int
	verbose     bool
)

var (
	gitConfigLocal, gitConfigRemote string
)

func init() {
	// connection
	flag.StringVar(&dbHost, "db-host", "localhost", "The ClickHouse server name. You can use either the name or the IPv4 or IPv6 address")
	flag.IntVar(&dbPort, "db-port", 9000, "The native ClickHouse port to connect to")
	flag.StringVar(&dbUsername, "db-username", "default", "The username. Default value: default")
	flag.StringVar(&dbPassword, "db-password", "", "The password. Default value: empty string")
	flag.BoolVar(&dbVerboseMode, "db-verbose", false, "Print query and other debugging info")

	// git
	flag.StringVar(&gitURL, "git-url", "", "URL of git repo with SQL queries")
	flag.StringVar(&gitBranch, "git-branch", "main", "Branch of git repo to use for SQL queries")
	flag.StringVar(&gitPath, "git-path", "", "Path within git repo to locate SQL queries")
	flag.StringVar(&workDir, "work-dir", "/tmp/sql-cd/", "Local path for repo with SQL queries")
	flag.StringVar(&gitPrivateKeyFile, "private-key-file", "/tmp/key", "Local path for the ssh private key")
	// flag.StringVar(&gitDest, "git-dest", "", "local path for repo with SQL queries")

	// config
	flag.StringVar(&gitConfigRemote, "remote-config", "", "Path to config")
	flag.StringVar(&gitConfigLocal, "local-config", "", "Path to config")

	// daemon
	flag.BoolVar(&runAsDaemon, "daemon", false, "Run as daemon")
	flag.IntVar(&timeout, "timeout", 60, "Global command timeout")
	flag.BoolVar(&verbose, "verbose", false, "Makes sql-cd verbose during the operation")

	flag.Parse()
}

func main() {
	var configMain config.Config

	for {
		if gitConfigLocal != "" && gitConfigRemote != "" {
			logging.Error.Fatalln("nope")
		}

		var configPath string

		if gitConfigLocal != "" || gitConfigRemote != "" {

			if gitConfigRemote != "" {
				if gitURL == "" {
					logging.Error.Fatalln("Please use -git-url flag")
				}
				// Clone config
				gitDest := getWorkDirName(workDir, gitURL, gitBranch)
				err := git.Clone(gitDest, gitURL, gitBranch, gitPrivateKeyFile, verbose)
				checkErr(err, runAsDaemon)

				configPath = filepath.Join(gitDest, gitConfigRemote)
			} else {
				_, err := os.Stat(gitConfigLocal)
				checkErr(err, runAsDaemon)
				configPath = gitConfigLocal
			}

			var err error
			logging.Debug.Println("config: ", configPath)
			configMain, err = config.ReadConfigFile(configPath)
			checkErr(err, runAsDaemon)

		} else {

			// Generate config from cli parameters
			configMain = config.Config{
				[]config.Cluster{
					config.Cluster{
						Name: dbHost,
						Host: dbHost,
						Port: dbPort,
						User: dbUsername,
						Pass: dbPassword,
						Sources: []config.Source{
							config.Source{
								GitRepo:   gitURL,
								GitBranch: gitBranch,
								GitPaths: []string{
									gitPath,
								},
							},
						},
					},
				},
			}
		}

		for _, cluster := range configMain.Clusters {
			logging.Debug.Println("Connect to ", cluster.Host)
			// Connect to ClickHouse
			ctx, conn, err := connect.Clickhouse(cluster.Host, cluster.Port, cluster.User, cluster.Pass, false)
			checkErr(err, runAsDaemon)
			for _, source := range cluster.Sources {
				gitDest := getWorkDirName(workDir, gitURL, gitBranch)
				err := git.Clone(gitDest, source.GitRepo, source.GitBranch, gitPrivateKeyFile, verbose)
				checkErr(err, runAsDaemon)
				for _, path := range source.GitPaths {
					logging.Debug.Println("Apply: ", path)
					// Apply SQL files
					queriesDir := filepath.Join(gitDest, path)
					err = apply.QueriesFromDir(ctx, conn, queriesDir, runAsDaemon)
					checkErr(err, runAsDaemon)
				}
				// Close connection
				conn.Close()
				logging.Info.Printf("%s connection closed", cluster.Host)
			}

		}

		// dev
		os.Exit(0)

		// Wait before next iteration
		if verbose {
			logging.Debug.Printf("Timeout %d sec\n", timeout)
		}
		time.Sleep(time.Duration(timeout) * time.Second)

	}
}

func getWorkDirName(subPaths ...string) string {

	var path string

	for _, subPath := range subPaths {
		path = filepath.Join(path, subPath)
	}

	replacer := strings.NewReplacer("https://", "", ":", "")

	return replacer.Replace(path)
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
