

```bash
goctl model mysql \
  datasource -url "user:pass@tcp(host:3306)/db" \
  -table "<table>" \
  -dir ./model \
  --style go_zero
```


```bash
goctl model mysql \
  ddl -src <file>.sql \
  -dir ./model \
  --style go_zero
```

**With cache:** add `-cache` flag to any command above.


```bash
goctl model mysql \
  ddl -c -src platform.sql \
  -dir .
```

