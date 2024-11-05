package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}

	cfg, err := config.ParseFlags()
	if err != nil {
		log.Fatalf(err.Error(), "event", "read config")
	}

	config.ReadConfig(cfg)
	if err != nil {
		log.Fatalf(err.Error(), "event", "read config")
	}

	logger.SetLogLevel(cfg.LogLevel)
	var handler *handlers.ServiceHandlers

	tickerSaveData := time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)

	if cfg.DBDSN != "" {
		database := dbInit(ctx, cfg)
		defer database.Close()
		handler = handlers.NewHandlers(database)

	} else {
		memStorage := storage.NewMemStorage(cfg.FileStoragePath)
		handler = handlers.NewHandlers(memStorage)
		if cfg.Restore {
			err := memStorage.RestoreFromFile()
			if err != nil {
				log.Error(err)
			}
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					err := memStorage.SaveToFile()
					if err != nil {
						log.Error(err)
					}
					log.Info("Closing program saved data")
					return
				case <-tickerSaveData.C:
					err := memStorage.SaveToFile()
					if err != nil {
						log.Error(err)
					}
				}
			}
		}()
	}

	r := router.NewRouter(cfg, handler)
	go func() {
		err := http.ListenAndServe(cfg.FlagRunAddr, r)
		if err != nil {
			log.Fatal(err)
		}
	}()
	<-termChan
	log.Info("Closing main program")
	cancel()

	wg.Wait()
}

func dbInit(ctx context.Context, cfg *config.ClientFlags) *dbstorage.PostgresStorage {
	// ps := "postgres://username:userpassword@localhost:5432/dbname"
	database := dbstorage.NewPostgresStorage(ctx, cfg.DBDSN)

	db, err := sql.Open("postgres", cfg.DBDSN)
	if err != nil {
		log.Error("Failed to open DB", "error", err)
	}
	defer db.Close()

	if err := goose.Up(db, "./internal/server/dbstorage/migrations"); err != nil {
		log.Error("Failed to run migrations", "error", err)
	}
	return database
}
