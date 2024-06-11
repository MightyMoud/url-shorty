package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dogmatiq/ferrite"
	"github.com/glebarez/sqlite"
	"github.com/joho/godotenv"
	"github.com/ms-mousa/url-shorty/middleware"
	"github.com/ms-mousa/url-shorty/models"
	"github.com/ms-mousa/url-shorty/services"
	"gorm.io/gorm"
)

var DatabaseLocation = ferrite.String("DB_CONNECTION", "Connection string for the database; this is the location of sql file - when using Sqlite").
	Required()

type ShortyRequest struct {
	Url string `json:"url"`
	Tag string `json:"tag"`
}

func main() {
	env := os.Getenv("SHORTY_ENV")
	if "" == env {
		env = "development"
	}
	godotenv.Load(".env" + "." + env)
	ferrite.Init()

	db, err := gorm.Open(sqlite.Open(DatabaseLocation.Value()), &gorm.Config{
		TranslateError: true,
	})
	if err != nil {
		panic("I failed my parents")
	}
	db.Exec("PRAGMA journal_mode=WAL; PRAGMA foreign_keys = ON; pragma synchronous = normal; pragma temp_store = memory;")
	ctx := context.WithValue(context.Background(), "db", db)

	db.AutoMigrate(&models.Entry{})

	router := http.NewServeMux()
	router.HandleFunc("GET /{shorty}", func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(ctx)

		entryToFind := models.Entry{Short: r.PathValue("shorty")}
		_, error := services.FindEntry(ctx, &entryToFind)
		if error != nil {
			panic(error)
		}
		http.Redirect(w, r, entryToFind.Url, http.StatusPermanentRedirect)
	})

	router.HandleFunc("POST /shorten/", func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(ctx)
		shortyReq := ShortyRequest{}

		err := json.NewDecoder(r.Body).Decode(&shortyReq)
		if err != nil {
			panic(err)
		}

		newEntry := models.Entry{Url: shortyReq.Url, Tag: shortyReq.Tag}
		shortUrl := services.ShortenUrl(ctx, newEntry.Url)
		newEntry.Short = shortUrl

		_, error := services.AddEntry(ctx, &newEntry)
		if errors.Is(error, gorm.ErrDuplicatedKey) {
			panic(error)
		}
		w.Write([]byte(fmt.Sprintln(newEntry)))
	})

	server := http.Server{
		Addr:    ":3000",
		Handler: middleware.LoggerMiddleware(router),
	}

	fmt.Println("Server listening on 3000")
	server.ListenAndServe()
}
