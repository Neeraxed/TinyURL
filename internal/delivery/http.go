package delivery

import (
	"io"
	"net/http"
	"time"
	"tinyurl/pkg/middleware"

	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type App struct {
	logic Logic
	log   *zap.Logger
}

type Logic interface {
	GetLink(id string) (string, error)
	AddLink(link string) (string, error)
}

func NewApp(logic Logic, log *zap.Logger) *App {
	return &App{
		logic: logic,
		log:   log,
	}
}

func (app *App) ApplyRoutes(o *middleware.Onion) *httprouter.Router {

	router := httprouter.New()
	router.GET("/:linkID", o.Apply(app.GetLinkHandler))
	router.POST("/", o.Apply(app.AddLinkHandler))

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

func (app *App) GetLinkHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	recievedID := ps.ByName("linkID")
	link, err := app.logic.GetLink(recievedID)
	if err != nil {
		app.log.Error("Failed to get link from database",
			zap.String("id", recievedID),
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return
	}

	http.Redirect(w, r, link, http.StatusSeeOther)
}

func (app *App) AddLinkHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		app.log.Error("Bad request body",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
		return
	}
	recievedLink := (string)(bytesBody)

	link, err := app.logic.GetLink(recievedLink)
	if err == nil {
		app.log.Info("Link already exists in database",
			zap.String("link", link),
			zap.Time("time", time.Now()),
		)
		w.Write([]byte(link))
		return
	}

	id, err := app.logic.AddLink(recievedLink)
	if err != nil {
		app.log.Error("Failed to add link to database",
			zap.String("link", recievedLink),
			zap.String("id", id),
			zap.Time("time", time.Now()),
		)
	}
	w.Write([]byte(id))
}
