package repository

import (
	"database/sql"

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

func (st *Storage) Init() error {

	st.config.connStr = "user=tester password=tester dbname=tester sslmode=disable"
	db, err := sql.Open("postgres", st.config.connStr)
	if err == nil {
		st.db = db
	}
	return err
}

func (st *Storage) Close() error {
	return st.db.Close()
}

func (st *Storage) AddToDB(id, recievedLink string) error {

	_, err := st.db.Exec("insert into links (short_link, long_link) values ($1, $2)",
		"/"+id, recievedLink)
	return err
}

func (st *Storage) GetFromDB(long_link string) (string, error) {

	var link string
	error := st.db.QueryRow("select long_link from links where short_link = $1", long_link).Scan(&link)
	return link, error
}
