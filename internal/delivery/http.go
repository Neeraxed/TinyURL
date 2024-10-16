package delivery

import (
	"io"
	"log"
	"net/http"
)

type App struct {
	logic Logic
}

type Logic interface {
	GetLink(link string) (string, error)
	AddLink(id, link string) error
	CreateId() (string, error)
}

func NewApp(logic Logic) *App {
	return &App{
		logic: logic,
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
	link, err := app.logic.GetLink(recievedLink)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, "/"+link, http.StatusSeeOther)
}

func (app *App) PostHandler(w http.ResponseWriter, r *http.Request) {

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

	id, err := app.logic.CreateId()
	if err != nil {
		log.Println(err)
		return
	}

	err = app.logic.AddLink(id, recievedLink)
	if err != nil {
		log.Println(err)
	}
}
