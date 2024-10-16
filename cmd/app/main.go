package main

import (
	"log"
	"net/http"

	"tinyurl/internal/repository"
	"tinyurl/internal/usecase"
	web "tinyurl/pkg"

	"github.com/julienschmidt/httprouter"
)

func main() {
	st := repository.NewStorage()
	st.Init()
	defer st.Close()

	uc := usecase.NewUsecase(st)
	app := web.NewApp(uc)

	router := httprouter.New()
	router.HandlerFunc("GET", "/", app.GetHandler)
	router.HandlerFunc("POST", "/", app.PostHandler)

	err := http.ListenAndServe(":3333", nil)
	log.Fatal(err)
}
