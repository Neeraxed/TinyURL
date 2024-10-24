package repository

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Storage struct {
	db *sql.DB
}

type Config struct {
	connStr string
}

func NewStorage() *Storage {
	return &Storage{}
}

func ReadConfig(log *zap.Logger) *Config {
	dbuser := os.Getenv("DBUSER")
	dbpass := os.Getenv("DBPASSWORD")

	return &Config{
		connStr: "user=" + dbuser + " password=" + dbpass + " dbname=tester sslmode=disable",
	}
}

func (st *Storage) Init(config *Config) error {
	db, err := sql.Open("postgres", config.connStr)
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
		id, recievedLink)
	return err
}

func (st *Storage) GetFromDB(long_link string) (string, error) {
	var link string
	error := st.db.QueryRow("select long_link from links where short_link = $1", long_link).Scan(&link)
	return link, error
}
