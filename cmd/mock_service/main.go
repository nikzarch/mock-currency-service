package main

import (
	"github.com/nikzarch/mock-currency-service/internal/db"
	"github.com/nikzarch/mock-currency-service/internal/valute"
	"net/http"
)

func main() {
	pool := db.GetPool()
	defer pool.Close()
	generator := valute.NewGenerator()
	repository := valute.NewPostgresRepository(pool)
	service := valute.NewService(repository, *generator)
	handler := valute.NewHandler(service)
	http.Handle("/XML_daily", handler)
	http.ListenAndServe(":8080", nil)
}
