package repository

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Storage struct {
	config Config
	db     *sql.DB
}
type Config struct {
	connStr string
}

func NewStorage() *Storage {
	return &Storage{}
}

func (st *Storage) Init() {

	db, err := sql.Open("postgres", st.config.connStr)
	if err != nil {
		panic(err)
	}
	st.db = db
}

func (st *Storage) Close() {
	st.db.Close()
}

func (st *Storage) AddToDB(id, recievedLink string) error {

	_, err := st.db.Exec("insert into links (short_link, long_link) values ($1, $2)",
		id, recievedLink)
	if err != nil {
		panic(err)
	}
	return err
}

func (st *Storage) GetFromDB(long_link string) (string, error) {

	var link string
	error := st.db.QueryRow("select short_link from links where long_link = $1", long_link).Scan(&link)
	if error != nil {
		log.Println(error)
	}
	return link, error
}
