package main

import (
	"log"
	"net/http"

	"tinyurl/internal/delivery"
	"tinyurl/internal/repository"
	"tinyurl/internal/usecase"

	"github.com/julienschmidt/httprouter"
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

	router := httprouter.New()
	router.HandlerFunc("GET", "/", app.GetHandler)
	router.HandlerFunc("POST", "/", app.PostHandler)

	err = http.ListenAndServe(":3333", nil)
	log.Fatal(err)
}
