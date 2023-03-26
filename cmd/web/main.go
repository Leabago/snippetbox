package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"snippetbox/pkg/models/mysql"
	"text/template"
	"time"

	"github.com/golangcollege/sessions"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	session       *sessions.Session
	infoLog       *log.Logger
	errorLog      *log.Logger
	debugLog      *log.Logger
	snippets      *mysql.SnippetModel
	users         *mysql.UserModel
	templateCache map[string]*template.Template

	authenticatedUserID string
	flash               string
	maxLength           int
	minLength           int
	maxEmailLength      int
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func getEnv(name string, errorLog *log.Logger) string {
	varEnv := os.Getenv(name)
	if varEnv == "" {
		ErrDuplicateEmail := fmt.Errorf("empty environment variable %s", name)
		errorLog.Fatal(ErrDuplicateEmail)
	}
	return varEnv
}

func main() {
	fmt.Println("start app")
	logFormat := log.Ldate | log.Ltime | log.Lshortfile
	infoLog := log.New(os.Stdout, "INFO\t", logFormat)
	debugLog := log.New(os.Stdout, "DEBUG\t", logFormat)
	errorLog := log.New(os.Stderr, "ERROR\t", logFormat)

	APP_PORT := getEnv("APP_PORT", errorLog)
	DB_NAME := getEnv("DB_NAME", errorLog)
	DB_HOST := getEnv("DB_HOST", errorLog)
	DB_WEB_USER := getEnv("DB_WEB_USER", errorLog)
	DB_WEB_PASSWORD := getEnv("DB_WEB_PASSWORD", errorLog)
	dsn := fmt.Sprintf("%s:%s@(%s)/%s?parseTime=true", DB_WEB_USER, DB_WEB_PASSWORD, DB_HOST, DB_NAME)
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour
	session.Secure = true
	session.SameSite = http.SameSiteStrictMode

	app := &application{
		session:       session,
		infoLog:       infoLog,
		debugLog:      debugLog,
		errorLog:      errorLog,
		snippets:      &mysql.SnippetModel{DB: db},
		users:         &mysql.UserModel{DB: db},
		templateCache: templateCache,

		authenticatedUserID: "authenticatedUserID",
		flash:               "flash",
		maxLength:           100,
		minLength:           10,
		maxEmailLength:      254,
	}

	srv := &http.Server{
		Addr:     APP_PORT,
		ErrorLog: errorLog,
		Handler:  app.routes(),

		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", APP_PORT)
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
