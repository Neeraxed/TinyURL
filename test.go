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

type Usecase struct {
	connStr  string
	database *sql.DB
}

func (usecase *Usecase) Init() {
	usecase.connStr = "user=tester password=tester dbname=tester sslmode=disable"

	db, err := sql.Open("postgres", usecase.connStr)
	if err != nil {
		panic(err)
	}
	usecase.database = db
	defer usecase.database.Close()
}

func (usecase *Usecase) indexGet(w http.ResponseWriter, r *http.Request) {

	bytesBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte("bad request body"))
		return
	}

	recievedLink := (string)(bytesBody)
	link, err := usecase.GetFromDB(recievedLink)
	if err != nil {
		log.Println(err)
		return
	}

	http.Redirect(w, r, link, http.StatusSeeOther)
}

func (usecase *Usecase) indexPost(w http.ResponseWriter, r *http.Request) {

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

	link, err := usecase.GetFromDB(recievedLink)
	if err == nil {
		log.Println("Existing link")
		w.Write([]byte(link))
		return
	}

	usecase.AddToDB(id, recievedLink)
}

type ConnectDB struct {
	usecase *Usecase
}

func (usecase *Usecase) AddToDB(id, recievedLink string) {

	_, err := usecase.database.Exec("insert into links (short_link, long_link) values ($1, $2)",
		id, recievedLink)
	if err != nil {
		panic(err)
	}
}

func (usecase *Usecase) GetFromDB(long_link string) (string, error) {

	var link string
	error := usecase.database.QueryRow("select short_link from links where long_link = $1", long_link).Scan(&link)
	if error != nil {
		log.Println(error)
	}
	return link, error
}

func main() {
	usecase := Usecase{}
	usecase.Init()

	router := httprouter.New()
	router.HandlerFunc("GET", "/", usecase.indexGet)
	router.HandlerFunc("POST", "/", usecase.indexPost)

	err := http.ListenAndServe(":3333", nil)
	log.Fatal(err)
}
