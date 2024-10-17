package main

import (
	"log"
	"net/http"

	"tinyurl/internal/delivery"
	"tinyurl/internal/repository"
	"tinyurl/internal/usecase"
)

func main() {
	st := repository.NewStorage()
	err := st.Init()
	if err != nil {
		log.Println("Failed to initialize storage: " + err.Error())
		return
	}
	defer st.Close()

	uc := usecase.NewUsecase(st)
	app := delivery.NewApp(uc)

	router := app.ApplyRouts()

	err = http.ListenAndServe(":3333", router)
	log.Fatal(err)
}
