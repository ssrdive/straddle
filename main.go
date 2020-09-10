package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ssrdive/basara/pkg/models/mysql"
)

type application struct {
	errorLog     *log.Logger
	infoLog      *log.Logger
	secret       []byte
	runtimeEnv   string
	clockworkAPI string
	api          *mysql.ApiModel
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "user:password@tcp(host)/database_name?parseTime=true", "MySQL data source name")
	secret := flag.String("secret", "straddle", "Secret key for generating jwts")
	runtimeEnv := flag.String("renv", "prod", "Runtime environment mode")
	clockworkAPI := flag.String("capi", "clockwork", "Clockwork SMS API key")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	defer db.Close()

	app := &application{
		errorLog:     errorLog,
		infoLog:      infoLog,
		secret:       []byte(*secret),
		runtimeEnv:   *runtimeEnv,
		clockworkAPI: *clockworkAPI,
		api:          &mysql.ApiModel{DB: db},
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, err
}
