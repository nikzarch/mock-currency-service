package main

import (
	"context"
	"fmt"
	"github.com/nikzarch/mock-currency-service/internal/db"
	"github.com/nikzarch/mock-currency-service/internal/valute"
	"time"
)

func main() {
	pool := db.GetPool()
	defer pool.Close()
	generator := valute.NewGenerator()
	repository := valute.NewPostgresRepository(pool)
	service := valute.NewService(repository, *generator)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, err := service.GetReportByDate(time.Now(), ctx)
	if err != nil {
		fmt.Println(err)
	}
}
