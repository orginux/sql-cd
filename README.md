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
  -timeout int
        Global command timeout (default 60)
  -work-dir string
        Local path for repo with SQL queries (default "/tmp/ufo/")
```

## Build
```bash
make build-go
make build-image
```

## Testing
### Create server and apply queries:
```bash
make test-server-up
```

### Check logs
```bash
docker logs -f sql-cd
```

### Check result:
```bash
docker attach clickhouse-client
```

```sql
select count() from system.tables;
```
