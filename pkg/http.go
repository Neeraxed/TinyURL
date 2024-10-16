package web

import (
	"io"
	"log"
	"net/http"

	"tinyurl/internal/usecase"

	"github.com/teris-io/shortid"
)

type App struct {
	uc *usecase.Usecase
}

func NewApp(uc *usecase.Usecase) *App {
	return &App{
		uc: uc,
	}
}
func (app *App) GetHandler(w http.ResponseWriter, r *http.Request) {

	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Bad request body"))
		return
	}

	recievedLink := (string)(bytesBody)
	link, err := app.uc.GetLink(recievedLink)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, link, http.StatusSeeOther)
}

func (app *App) PostHandler(w http.ResponseWriter, r *http.Request) {

	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Bad request body"))
		return
	}
	recievedLink := (string)(bytesBody)

	id, ok := shortid.Generate()
	if ok != nil {
		log.Println("Did not generate")
	}
	id = "/" + id

	link, err := app.uc.GetLink(recievedLink)
	if err == nil {
		log.Println("Existing link")
		w.Write([]byte(link))
		return
	}

	app.uc.AddLink(id, recievedLink)
}
