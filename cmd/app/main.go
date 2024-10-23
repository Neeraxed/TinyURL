package main

import (
	"net/http"

	"tinyurl/internal/delivery"
	"tinyurl/internal/repository"
	"tinyurl/internal/usecase"
	"tinyurl/pkg/middleware"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	st := repository.NewStorage()
	err := st.Init()
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

	err = http.ListenAndServe(":3333", router)
	logger.Fatal("Server died",
		zap.String("message", err.Error()),
	)
}
