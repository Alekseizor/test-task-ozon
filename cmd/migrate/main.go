package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	jww "github.com/spf13/jwalterweatherman"
)

const (
	dbName         = "generation"
	dbUser         = "root"
	dbPass         = "ozon"
	dbHost         = "db"
	dbPort         = "5432"
	migrationsPath = "migrations"
	driver         = "postgres"
)

func main() {
	jww.SetLogThreshold(jww.LevelInfo)
	jww.SetStdoutThreshold(jww.LevelInfo)

	ctx := context.Background()
	jww.INFO.Println("Starting migrations")

	// подключаемся к БД
	db, err := connect()
	if err != nil {
		jww.ERROR.Fatalln(err)
	}
	jww.INFO.Println("The database connection was established successfully")

	// устанавливаем свой логер
	goose.SetLogger(&gooseLogger{ctx: ctx})
	// запускаем миграции
	jww.INFO.Println("Upping migrations")
	err = goose.SetDialect(driver)
	if err != nil {
		jww.ERROR.Fatalf("Failed to set dialect: %v", err)
	}
	err = goose.Up(db, migrationsPath)
	if err != nil {
		jww.ERROR.Fatalf("Failed to migrate: %v", err)
	}

	jww.INFO.Println("DB migration completed")
}

// Выполняет подключение к БД
func connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", dbUser, dbName, dbPass, dbHost, dbPort)
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Реализация интерфйса goose.Logger
type gooseLogger struct {
	ctx context.Context
}

func (gl *gooseLogger) Fatal(v ...interface{}) {
	jww.FATAL.Fatal(v...)
}
func (gl *gooseLogger) Fatalf(format string, v ...interface{}) {
	jww.FATAL.Fatalf(format, v...)
}
func (gl *gooseLogger) Print(v ...interface{}) {
	jww.INFO.Print(v...)
}
func (gl *gooseLogger) Println(v ...interface{}) {
	jww.INFO.Println(v...)
}
func (gl *gooseLogger) Printf(format string, v ...interface{}) {
	jww.INFO.Printf(format, v...)
}
