package main

import (
	"net/http"
	"os"

	"tinyurl/internal/delivery"
	"tinyurl/internal/repository"
	"tinyurl/internal/usecase"
	"tinyurl/pkg/middleware"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"

	"go.uber.org/zap"
)

func main() {
	godotenv.Load("./.env")
	logger, _ := zap.NewDevelopment()

	config := repository.ReadConfig(logger)
	st := repository.NewStorage()
	err := st.Init(config)
	if err != nil {
		logger.Fatal("Failed to initialize storage",
			zap.String("message", err.Error()),
		)
		return
	}
	defer st.Close()

	uc := usecase.NewUsecase(st, logger)
	app := delivery.NewApp(uc, logger)

	onion := middleware.NewOnion(logger)

	onion.AppendMiddleware(
		onion.Timer,
		onion.LogRequestResponse)

	router := app.ApplyRoutes(onion)

	port := os.Getenv("PORT")
	err = http.ListenAndServe(port, router)
	logger.Fatal("Server died",
		zap.String("message", err.Error()),
	)
}
