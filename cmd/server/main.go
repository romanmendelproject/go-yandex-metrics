package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/config"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/dbstorage"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/handlers"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/middlewares/logger"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/router"
	"github.com/romanmendelproject/go-yandex-metrics/internal/server/storage"

	_ "github.com/romanmendelproject/go-yandex-metrics/internal/server/dbstorage/migrations"

	log "github.com/sirupsen/logrus"
)

var buildVersion string
var buildDate string
var buildCommit string

func printVersion() {
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func main() {
	printVersion()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.ParseFlags()
	logger.SetLogLevel(config.LogLevel)
	var handler *handlers.ServiceHandlers

	if config.DBDSN != "" {
		database := dbInit(ctx)
		defer database.Close()
		handler = handlers.NewHandlers(database)

	} else {
		memStorage := storage.NewMemStorage(config.FileStoragePath)
		handler = handlers.NewHandlers(memStorage)
		if config.Restore {
			err := memStorage.RestoreFromFile()
			if err != nil {
				log.Error(err)
			}
		}
		go func() {
			for {
				time.Sleep(time.Second * time.Duration(config.StoreInterval))
				err := memStorage.SaveToFile()
				if err != nil {
					log.Error(err)
				}
			}
		}()

	}
	r := router.NewRouter(handler)
	func() {

		err := http.ListenAndServe(config.FlagRunAddr, r)
		if err != nil {
			log.Fatal(err)
		}
	}()
}

func dbInit(ctx context.Context) *dbstorage.PostgresStorage {
	// ps := "postgres://username:userpassword@localhost:5432/dbname"

	database := dbstorage.NewPostgresStorage(ctx, config.DBDSN)

	db, err := sql.Open("postgres", config.DBDSN)
	if err != nil {
		log.Error("Failed to open DB", "error", err)
	}
	defer db.Close()

	if err := goose.Up(db, "./internal/server/dbstorage/migrations"); err != nil {
		log.Error("Failed to run migrations", "error", err)
	}
	return database
}
