package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
	"net/http"

	"skud/app"
	"skud/internal/config"
	sqlRepo "skud/internal/repository/sql"
	"skud/service"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}
	db, err := sqlx.Open("mysql", cfg.DBUrl())
	if err != nil {
		log.Fatalf("failed to open DB: %s", err)
	}
	repo := sqlRepo.New(db)
	svc := service.New(repo)
	router := app.NewHTTPRouter(svc)
	log.Printf("listening on :%s", cfg.HTTPPort)
	if err = http.ListenAndServe(":"+cfg.HTTPPort, router); err != nil {
		log.Fatal(err)
	}
}
