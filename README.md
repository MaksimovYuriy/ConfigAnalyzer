# ConfigAnalyzer

Проверка YAML- и JSON-конфигураций через CLI, REST и gRPC.
Утилита находит debug-логирование, открытые пароли, отключённый TLS и проверку
сертификатов, небезопасные алгоритмы, прослушивание `0.0.0.0` и избыточные права файла.
Для каждой проблемы выводятся путь, уровень LOW/MEDIUM/HIGH и рекомендация.

## Сборка

Нужен Go 1.27.3. Команды выполняются из корня проекта в Bash.

```sh
go mod download
mkdir -p bin
go build -o bin/config-analyzer-cli ./cmd/analyzer
go build -o bin/config-analyzer-app ./cmd/app
```

Сгенерированный Protobuf-код включён в проект: для обычной сборки `protoc` не нужен.

## CLI

```sh
./bin/config-analyzer-cli examples/safe.yaml
./bin/config-analyzer-cli examples/insecure.yaml
./bin/config-analyzer-cli examples/insecure.yaml --silent
./bin/config-analyzer-cli --stdin < examples/debug.yaml
```

Для ручного ввода запустите `--stdin`, вставьте конфиг и нажмите Enter, затем Ctrl+D.
Формат определяется по содержимому, расширение файла не используется.

| Флаг | Назначение |
| --- | --- |
| `--stdin` | Читать конфигурацию из стандартного ввода вместо файла |
| `-s`, `--silent` | Сохранить отчёт, но не возвращать ошибку из-за находок |
| `-h`, `--help` | Справка |

Путь — первый позиционный аргумент. Его нельзя совмещать с `--stdin`.
Коды завершения: **0** — нет находок или указан `--silent`, **1** — есть находки,
**2** — ошибка аргументов, чтения или разбора; `--silent` такие ошибки не подавляет.

`examples/safe.yaml` выводит `Анализ завершён. Проблем не найдено.` и код `0`.
`examples/insecure.yaml` даёт шесть находок и код `1`.
`examples/debug.yaml` даёт одну находку:

```text
LOW: services[0].logging.level: Логирование в debug-режиме. Поменяйте уровень логирования на info или выше.
```

При передаче локального пути дополнительно проверяются права через `os.Stat`:
запись для остальных пользователей (`o+w`, например `0666`) даёт HIGH.
Для stdin, REST и gRPC эта проверка пропускается.

## REST и gRPC

```sh
./bin/config-analyzer-app
```

Один процесс запускает оба сервера. Ctrl+C/SIGTERM корректно останавливает их;
по истечении настроенного таймаута оставшиеся запросы прерываются.

**REST:** `POST http://127.0.0.1:8081/analyze`, тело — form-data с текстовым полем `content`.

```sh
curl http://127.0.0.1:8081/analyze \
  --form 'content=<examples/insecure.yaml'
```

В Postman: **Body → form-data → content → Text**, затем вставьте конфиг в Value.
Content-Type с boundary задаётся автоматически. Ответ — JSON-массив проблем,
без находок — `[]`. Неверный конфиг даёт `400`, превышение лимита тела 1 MiB — `413`.

**gRPC:** адрес `127.0.0.1:50051`, метод `grpc.v1.ConfigService/Analyze`.
В Postman выберите gRPC-запрос без TLS; описание доступно через reflection
или импорт `docs/proto/v1/analyze.proto`. Сообщение:

```json
{
  "content": "logging:\n  level: debug"
}
```

Ответ содержит список `problems`. Неверный конфиг даёт `InvalidArgument`,
превышение лимита сообщения — `ResourceExhausted`.
Оба сервера по умолчанию работают без TLS.

## Настройки серверов

Настройки читаются через `cleanenv` из окружения и необязательного `.env`
в рабочем каталоге. Значения из `.env` перезаписывают одноимённые переменные окружения.

| Переменная | По умолчанию |
| --- | --- |
| `HTTP_ADDRESS` / `HTTP_PORT` | `0.0.0.0` / `8081` |
| `GRPC_ADDRESS` / `GRPC_PORT` | `0.0.0.0` / `50051` |
| `HTTP_READ_TIMEOUT` / `HTTP_WRITE_TIMEOUT` | `60s` / `60s` |
| `HTTP_READ_HEADER_TIMEOUT` / `HTTP_IDLE_TIMEOUT` | `5s` / `60s` |
| `HTTP_SHUTDOWN_TIMEOUT` / `GRPC_SHUTDOWN_TIMEOUT` | `10s` / `10s` |
| `GRPC_MAX_RECV_MSG_SIZE` | `1048576` байт |

## Проверка проекта

```sh
go vet ./...
go test -race ./...
bash scripts/smoke.sh ./bin/config-analyzer-cli
```

Для race detector нужен C-компилятор. Smoke-скрипт проверяет запуск CLI и коды
завершения; успешный результат — `Smoke checks passed.`.
[CI](.github/workflows/ci.yml) выполняет проверку зависимостей, форматирования,
`go vet`, тесты с race detector, сборку CLI и приложения, затем smoke-проверки.
Запускается при push, pull request и вручную во вкладке Actions.

## Допущения

- Имена параметров и алгоритмы проверяются по спискам в `internal/rules`.
- Ссылки `${VAR}`, `{{ secret }}` и теги `!vault`, `!env`, `!secret` не считаются
  открытыми паролями. Значения паролей в отчёт не выводятся.
- `0.0.0.0` даёт MEDIUM: внешние сетевые ограничения по конфигу не определяются.
- Поддерживаются вложенные объекты, массивы, скалярный корень и YAML-алиасы;
  YAML merge (`<<`) не разворачивается. Один ввод — один документ.
- Проверка прав учитывает Unix-биты, без ACL и прав родительских каталогов.
