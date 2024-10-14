package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/lib/pq" //To import a package solely for its side-effects (initialization), use the blank identifier as explicit package name - хз вообще что значит
	"github.com/teris-io/shortid"
)

var database *sql.DB

func main() {
	//Database connection
	connStr := "user=tester password=tester dbname=tester sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	database = db // хз зачем (почему database, err := sql.Open(...) не работало), надо разобраться
	defer db.Close()

	//HTTP server
	http.HandleFunc("/", handler)
	err1 := http.ListenAndServe(":3333", nil)
	if err1 != nil {
		fmt.Println("ашипка: ", err1)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {

		bytesBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			w.Write([]byte("bad request body"))
			return
		}
		recievedLink := (string)(bytesBody)

		//Check if link already exists
		var link string
		err = database.QueryRow("select short_link from links where long_link = $1", recievedLink).Scan(&link)
		if err == nil {
			log.Println("existing link")
			w.Write([]byte(link))
			return
		}

		id, ok := shortid.Generate()
		if ok != nil {
			fmt.Println("did not generate")
		}
		id = "/" + id

		w.Write([]byte(id))

		//Database write
		_, err = database.Exec("insert into links (short_link, long_link) values ($1, $2)",
			id, recievedLink)
		if err != nil {
			panic(err)
		}

	} else if r.Method == "GET" {

		//Database read
		var link string
		err := database.QueryRow("select long_link from links where short_link = $1", r.URL.Path).Scan(&link)
		if err != nil {
			log.Println(err)
		}

		http.Redirect(w, r, link, http.StatusSeeOther)

	} else {

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("неправильный метод запроса"))
	}
}
