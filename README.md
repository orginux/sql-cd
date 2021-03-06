# SQL as a code
[![Go](https://github.com/orginux/sql-cd/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/orginux/sql-cd/actions/workflows/go.yml)


## Usage
```bash
  -daemon
        Run as daemon
  -db-host string
        The ClickHouse server name. You can use either the name or the IPv4 or IPv6 address (default "localhost")
  -db-password string
        The password. Default value: empty string
  -db-port int
        The native ClickHouse port to connect to (default 9000)
  -db-username string
        The username. Default value: default (default "default")
  -db-verbose
        Print query and other debugging info
  -git-branch string
        Branch of git repo to use for SQL queries (default "main")
  -git-path string
        Path within git repo to locate SQL queries
  -git-url string
        URL of git repo with SQL queries
  -private-key-file string
        Local path for the ssh private key (default "/tmp/key")
  -timeout int
        Global command timeout (default 60)
  -verbose
        Makes sql-cd verbose during the operation
  -work-dir string
        Local path for repo with SQL queries (default "/tmp/sql-cd/")
```

## Build
```bash
make build-go
make build-image
```

## Testing
Check the `tests/` directory.
