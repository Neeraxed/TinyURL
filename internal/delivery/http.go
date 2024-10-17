package delivery

import (
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type App struct {
	logic Logic
}

type Logic interface {
	GetLink(link string) (string, error)
	AddLink(link string) (string, error)
}

func NewApp(logic Logic) *App {
	return &App{
		logic: logic,
	}
}

func (app *App) ApplyRouts() *httprouter.Router {
	router := httprouter.New()
	router.GET("/", app.GetLinkHandler)
	router.POST("/", app.AddLinkHandler)

	router.GlobalOPTIONS = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Access-Control-Request-Method") != "" {
			// Set CORS headers
			header := w.Header()
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
			header.Set("Access-Control-Allow-Origin", "*")
		}

		// Adjust status code to 204
		w.WriteHeader(http.StatusNoContent)
	})

	return router
}

func (app *App) GetLinkHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Bad request body"))
		return
	}

	recievedLink := (string)(bytesBody)
	link, err := app.logic.GetLink(recievedLink)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, link, http.StatusSeeOther)
}

func (app *App) AddLinkHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Bad request body"))
		return
	}
	recievedLink := (string)(bytesBody)

	link, err := app.logic.GetLink(recievedLink)
	if err == nil {
		log.Println("Existing link")
		w.Write([]byte(link))
		return
	}

	id, err := app.logic.AddLink(recievedLink)
	if err != nil {
		log.Println(err)
	}

	w.Write([]byte(id))
}
