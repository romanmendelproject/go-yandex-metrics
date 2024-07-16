package main

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"

	"github.com/romanmendelproject/go-yandex-metrics/internal/server/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/dbstorage"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/handlers"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/logger"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/router"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/storage"

	_ "github.com/romanmendelproject/go-yandex-metrics/internal/server/migrations"

	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.ParseFlags()
	logger.SetLogLevel(config.LogLevel)
	var handler *handlers.ServiceHandlers

	if config.DBDSN != "" {
		// ps := "postgres://username:userpassword@localhost:5432/dbname"

		database := dbstorage.NewDBStorage(ctx, config.DBDSN)
		defer database.Close()

		db, err := sql.Open("postgres", config.DBDSN)
		if err != nil {
			log.Error("Failed to open DB", "error", err)
		}
		defer db.Close()

		if err := goose.Up(db, "./internal/server/migrations"); err != nil {
			log.Error("Failed to run migrations", "error", err)
		}

		handler = handlers.NewHandlers(database)
		runServer(handler)
	} else {
		memStorage := storage.NewMemStorage(config.FileStoragePath)
		handler = handlers.NewHandlers(memStorage)
		if config.Restore {
			err := memStorage.RestoreFromFile()
			if err != nil {
				log.Error(err)
			}
		}
		runServer(handler)

		for {
			time.Sleep(time.Second * time.Duration(config.StoreInterval))
			go func() {
				err := memStorage.SaveToFile()
				if err != nil {
					log.Error(err)
				}
			}()
		}
	}
}

func runServer(handler *handlers.ServiceHandlers) {
	r := router.NewRouter(handler)
	func() {

		err := http.ListenAndServe(config.FlagRunAddr, r)
		if err != nil {
			panic(err)
		}
	}()
}
