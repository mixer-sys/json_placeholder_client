# json_placeholder_client
Клиент для JSONPlaceholder

## Start
Для быстрого старта выполните команды
```bash
mv env.example .env
docker build -t proxy .
docker run -p 8080:8080 --env-file .env proxy
```