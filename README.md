# ConfigAnalyzer

Утилита на Go для поиска небезопасных настроек в YAML и JSON.

Проект находится на этапе подготовки структуры: анализ конфигураций ещё не
реализован. Сейчас доступна справка; остальные вызовы завершаются с кодом `2`.

## Запуск

Нужна версия Go, указанная в `go.mod` (также закреплена в `.tool-versions`).

```sh
go run ./cmd/config-analyzer --help
```

Сборка:

```sh
mkdir -p bin
go build -o bin/config-analyzer ./cmd/config-analyzer
```

## Структура

- `cmd/config-analyzer` — точка входа приложения.
- `internal/cli` — аргументы, вывод и коды завершения CLI.
- `internal/config` — чтение YAML и JSON.
- `internal/analyzer` — модель проблемы и общий анализатор.
- `internal/rules` — правила проверки и варианты названий параметров.
- `.github/workflows/ci.yml` — форматирование, vet, тесты и сборка.

## Проверка

```sh
go vet ./...
go test ./...
go build -o bin/config-analyzer ./cmd/config-analyzer
```

Тестов пока нет. Они будут добавлены вместе с реализацией поведения.
