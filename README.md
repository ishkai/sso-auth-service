# SSO Auth Service

SSO Auth Service - сервис авторизации на Go, который предоставляет gRPC API для регистрации пользователей, входа в систему, генерации JWT-токенов и проверки прав администратора.

## Возможности

* Регистрация пользователей
* Авторизация пользователей
* Хеширование паролей через bcrypt
* Генерация JWT access token
* Поддержка app-based JWT secret
* Проверка admin-статуса пользователя
* PostgreSQL storage
* SQLite storage
* SQL-миграции
* gRPC API
* gRPC health check
* gRPC reflection
* Prometheus metrics
* Docker Compose окружение
* Graceful shutdown
* Структурированное логирование через `slog`
* Интеграционные тесты

## Стек

* Go
* gRPC
* Protocol Buffers
* PostgreSQL
* SQLite
* JWT
* bcrypt
* Docker
* Docker Compose
* Prometheus
* golang-migrate
* slog

## Архитектура

```text
Client
  |
  | gRPC
  v
SSO Auth Service
  |
  | SQL
  v
PostgreSQL
```

При запуске через Docker Compose поднимаются следующие сервисы:

```text
postgres       - база данных PostgreSQL
migrator       - контейнер для применения миграций
auth-service   - gRPC сервис авторизации
prometheus     - сбор метрик auth-service
```

## Структура проекта

```text
.
├── cmd/
│   ├── sso/              # точка входа auth-service
│   └── migator/          # запуск миграций
├── config/               # конфигурационные файлы
├── deploy/
│   └── prometheus/       # конфигурация Prometheus
├── internal/
│   ├── app/              # инициализация приложения
│   ├── config/           # загрузка конфига
│   ├── domain/           # доменные модели
│   ├── grpc/             # gRPC handlers
│   ├── lib/              # общие библиотеки
│   ├── observability/    # метрики и мониторинг
│   ├── services/         # бизнес-логика
│   └── storage/          # реализации хранилища
├── migrations/           # SQL-миграции
├── tests/                # интеграционные тесты
├── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## Конфигурация

Пример конфигурации для Docker:

```yaml
env: "local"

storage_type: postgres
storage_dsn: "${STORAGE_DSN}"

token_ttl: 1h

grpc:
  port: $GRPC_PORT
  timeout: 5s

metrics:
  port: $METRICS_PORT
```


## Запуск через Docker Compose

```bash
docker compose up --build
```


Посмотреть логи auth-service:

```bash
docker compose logs  auth-service
```

Посмотреть логи migrator:

```bash
docker compose logs migrator
```

Остановить сервисы:

```bash
docker compose down
```

## Миграции

В Docker Compose миграции применяются автоматически через контейнер `migrator`.

Ручной запуск миграций:

```bash
go run ./cmd/migator \
  --database-url "$STORAGE_DSN" \
  --migrations-path ./migrations/postgres
```

## gRPC Reflection

Сервис поддерживает gRPC reflection. Благодаря этому API можно смотреть и вызывать через `grpcurl` без локальных `.proto` файлов.

Показать список сервисов:

```bash
grpcurl -plaintext localhost:$GRPC_PORT list
```

Ожидаемый результат:

```text
auth.Auth
grpc.health.v1.Health
grpc.reflection.v1.ServerReflection
```

Показать методы auth-сервиса:

```bash
grpcurl -plaintext localhost:$GRPC_PORT list auth.Auth
```

Показать описание auth-сервиса:

```bash
grpcurl -plaintext localhost:$GRPC_PORT describe auth.Auth
```

## gRPC Health Check

Проверить общий статус gRPC-сервера:

```bash
grpcurl -plaintext \
  -d '{"service":""}' \
  localhost:$GRPC_PORT \
  grpc.health.v1.Health/Check
```

Ожидаемый ответ:

```json
{
  "status": "SERVING"
}
```

Проверить статус auth-сервиса:

```bash
grpcurl -plaintext \
  -d '{"service":"auth.Auth"}' \
  localhost:$GRPC_PORT \
  grpc.health.v1.Health/Check
```

Ожидаемый ответ:

```json
{
  "status": "SERVING"
}
```

## Примеры API

### Register

```bash
grpcurl -plaintext \
  -d '{"email":"user@example.com","password":"123456"}' \
  localhost:$GRPC_PORT \
  auth.Auth/Register
```

Пример ответа:

```json
{
  "userId": "1"
}
```

### Login

```bash
grpcurl -plaintext \
  -d '{"email":"user@example.com","password":"123456","appId":1}' \
  localhost:$GRPC_PORT \
  auth.Auth/Login
```

Пример ответа:

```json
{
  "token": "jwt-token"
}
```

### IsAdmin

```bash
grpcurl -plaintext \
  -d '{"userId":1}' \
  localhost:$GRPC_PORT \
  auth.Auth/IsAdmin
```

Пример ответа:

```json
{
  "isAdmin": false
}
```

## Метрики

Сервис отдает Prometheus-метрики по адресу:

```text
http://localhost:$METRICS_PORT/metrics
```

Проверить метрики вручную:

```bash
curl localhost:$METRICS_PORT/metrics
```

Показать только метрики сервиса:

```bash
curl -s localhost:$METRICS_PORT/metrics | grep sso
```

Prometheus доступен по адресу:

```text
http://localhost:$PROMETHEUS_PORT
```

Примеры PromQL-запросов:

```promql
sso_auth_events_total
```

```promql
sso_grpc_requests_total
```

```promql
rate(sso_grpc_requests_total[1m])
```

```promql
histogram_quantile(
  0.95,
  sum(rate(sso_grpc_request_duration_seconds_bucket[5m])) by (le, method)
)
```

## Основные метрики

```text
sso_grpc_requests_total
sso_grpc_request_duration_seconds
sso_auth_events_total
```

Сервис собирает:

* количество gRPC-запросов
* длительность gRPC-запросов
* коды ответов gRPC
* события регистрации
* события логина
* успешные и неуспешные auth-операции

## Тесты

Запустить все тесты:

```bash
go test ./...
```







