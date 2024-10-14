package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq" //To import a package solely for its side-effects (initialization), use the blank identifier as explicit package name - хз вообще что значит
	"github.com/teris-io/shortid"
)

type App struct {
	config   Config
	database *sql.DB
}
type Config struct {
	connStr string
}

func (app *App) Init() {
	app.config.connStr = "user=tester password=tester dbname=tester sslmode=disable"

	db, err := sql.Open("postgres", app.config.connStr)
	if err != nil {
		panic(err)
	}
	app.database = db
}

func (app *App) Close() {
	app.database.Close()
}

func (uc *Usecase) GetHandler(w http.ResponseWriter, r *http.Request) {

	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Bad request body"))
		return
	}

	recievedLink := (string)(bytesBody)
	link, err := uc.GetFromDB(recievedLink)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, link, http.StatusSeeOther)
}

func (uc *Usecase) PostHandler(w http.ResponseWriter, r *http.Request) {

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

	link, err := uc.GetFromDB(recievedLink)
	if err == nil {
		log.Println("Existing link")
		w.Write([]byte(link))
		return
	}

	uc.AddToDB(id, recievedLink)
}

type Usecase struct {
	app *App
}

func (uc *Usecase) AddToDB(id, recievedLink string) {

	_, err := uc.app.database.Exec("insert into links (short_link, long_link) values ($1, $2)",
		id, recievedLink)
	if err != nil {
		panic(err)
	}
}

func (uc *Usecase) GetFromDB(long_link string) (string, error) {

	var link string
	error := uc.app.database.QueryRow("select short_link from links where long_link = $1", long_link).Scan(&link)
	if error != nil {
		log.Println(error)
	}
	return link, error
}

func main() {
	app := App{}
	app.Init()
	defer app.Close()

	uc := Usecase{} //сомнительно

	router := httprouter.New()
	router.HandlerFunc("GET", "/", uc.GetHandler)
	router.HandlerFunc("POST", "/", uc.PostHandler)

	err := http.ListenAndServe(":3333", nil)
	log.Fatal(err)
}
