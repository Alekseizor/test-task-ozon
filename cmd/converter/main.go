package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"test-task-ozon/internal/pkg/handlers"
	"test-task-ozon/internal/pkg/repository/links"
	"test-task-ozon/internal/pkg/sendingjson"
)

const (
	//По-хорошему это все в .env файл надо, но пока для наглядности здесь
	dbName = "generation"
	dbUser = "root"
	dbPass = "ozon"
	dbHost = "db"
	dbPort = "5432"
	driver = "postgres"
)

func main() {
	ctx := context.Background()
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", dbUser, dbName, dbPass, dbHost, dbPort)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Println(fmt.Errorf("failed to connect to the db - %s", err.Error()))
		return
	}
	db.SetMaxOpenConns(0)
	err = db.Ping()
	if err != nil {
		log.Println(fmt.Errorf("failed to connect to the db - %s", err.Error()))
		return
	}
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Println(fmt.Errorf("couldn't create a new logger - %s", err.Error()))
		return
	}
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()
	serviceSend := sendingjson.NewServiceSend(logger)
	linkHandler := &handlers.LinksHandler{
		Logger: logger,
		Send:   serviceSend,
	}

	value := os.Getenv("METHOD")
	if value == "in-memory" {
		repoLinkInMemory, err := links.NewRepoLinkInMemory()
		if err != nil {
			logger.Errorf("failed to create NewRepoLinkInMemory - %s", err.Error())
		}
		linkHandler.LinkRepo = repoLinkInMemory
		go handlers.StartConverterServer(repoLinkInMemory)
	} else {
		repoLinkPostgres, err := links.NewRepoLinkPostgres(db, ctx)
		if err != nil {
			logger.Errorf("failed to create NewRepoLinkPostgres - %s", err.Error())
		}
		linkHandler.LinkRepo = repoLinkPostgres
		go handlers.StartConverterServer(repoLinkPostgres)
	}

	r := mux.NewRouter()

	r.HandleFunc("/{URL}", linkHandler.GetLink).Methods("GET")
	r.HandleFunc("/api/links", linkHandler.Generation).Methods("POST")

	addr := ":8080"
	logger.Infow("starting server",
		"type", "START",
		"addr", addr,
	)
	err = http.ListenAndServe(addr, r)
	if err != nil {
		logger.Errorf("couldn't start listening - %s", err.Error())
	}
}
