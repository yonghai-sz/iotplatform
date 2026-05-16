
```sh
bash scripts/dev/logs.sh internal-api
bash scripts/dev/logs.sh emqx
bash scripts/dev/logs.sh mqtt-subscriber
bash scripts/dev/logs.sh internal-api emqx mqtt-subscriber
```

<https://docs.emqx.com/en/emqx/v5.8/access-control/authn/http.html>



```sh
docker compose -f scripts/dev/docker-compose.yml up -d --build mqtt-subscriber
```