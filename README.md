# go-musthave-metrics-tpl

Репозитория для трека «Сервер сбора метрик и алертинга».

## Запуск программы
1. Клонируем репозиторий и переходим в него
2. Запускаем БД
- cd docker
- docker-compose up -d
3. Запуск сервера
- go run cmd/server/main.go
4. Запуск агента
- go run cmd/agent/main.go

## Запуск тестов
1. Клонируем репозиторий и переходим в него
2. Запускаем БД
- cd docker
- docker-compose up -d
3. go test ./...  -coverprofile cover.out 
4. go tool cover -func cover.out

## Документация по анализаторам
1. Клонируем репозиторий и переходим в него
2. godoc -http=:8080
3. http://127.0.0.1:8080/pkg/github.com/romanmendelproject/go-yandex-metrics/cmd/staticlint/
