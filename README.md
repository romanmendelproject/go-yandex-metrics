# go-musthave-metrics-tpl

Репозитория для трека «Сервер сбора метрик и алертинга».

## Запуск тестов

1. Клонируем репозиторий и ререходим в него
2. go test ./...  -coverprofile cover.out 
3. go tool cover -func cover.out