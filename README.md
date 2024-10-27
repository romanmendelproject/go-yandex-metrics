# go-musthave-metrics-tpl

Репозитория для трека «Сервер сбора метрик и алертинга».

## Запуск тестов

1. Клонируем репозиторий и переходим в него
2. Запускаем БД
- cd docker
- docker-compose up -d
3. go test ./...  -coverprofile cover.out 
4. go tool cover -func cover.out