# Run with parameters from flags
This is a simple sql-cd use case, we are passing parameters via cli flags. In this example, we are downloading [the repo](https://github.com/orginux/sql-cd-test) over an HTTPS connection.

Run test enviroment:
```bash
docker-compose up -d --build
```

You can check logs:
```bash
docker logs -f sql-cd
```

And check results, for example:
```bash
docker attach clickhouse-client
```

```bash
select name from system.tables where database = 'common_db';
```
