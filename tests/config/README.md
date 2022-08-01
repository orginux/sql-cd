# Remote config fale test case
In this case we start 2 ClickHouse clusters version 22.3.
Then sql-cd downloads the configuration file from [the repo](https://github.com/orginux/sql-cd-test) and tries to apply it. In this example, we are downloading the repo via an SSH connection and key are requirement, this command helps you put your key into container:

```bash
docker cp github-read sql-cd:/tmp/key
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
