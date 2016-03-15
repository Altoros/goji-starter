package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-community/go-cfenv"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/net/context"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

const (
	DEFAULT_PORT = 8080
)

type comment struct {
	ID     int64  `json:"id"`
	Author string `json:"author"`
	Text   string `json:"text"`
}

var index = template.Must(template.ParseFiles(
	"templates/_base.html",
	"templates/index.html",
))

func handleIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	index.Execute(w, nil)
}

func readCommentsFromDb(w http.ResponseWriter) {
	comments, err := db.GetComments()
	if err != nil {
		log.Printf("Unable to read comments - %s", err)
		http.Error(w, fmt.Sprintf("Unable to read comments - %s", err), http.StatusInternalServerError)
		return
	}
	commentData, err := json.MarshalIndent(comments, "", "    ")
	if err != nil {
		log.Printf("Unable to marshal comments to json: %s", err)
		http.Error(w, fmt.Sprintf("Unable to marshal comments to json: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	io.Copy(w, bytes.NewReader(commentData))
}

func getComments(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	readCommentsFromDb(w)
}

func postComments(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	err := db.AddComment(comment{ID: time.Now().UnixNano() / 1000000, Author: r.FormValue("author"), Text: r.FormValue("text")})
	if err != nil {
		log.Printf("Unable to store comment - %s", err)
		http.Error(w, fmt.Sprintf("Unable to store comment - %s", err), http.StatusInternalServerError)
		return
	}
	readCommentsFromDb(w)
}

var db *DbConnection

func main() {
	var connString string
	var port int

	if appEnv, err := cfenv.Current(); err != nil {
		connString = "sslmode=disable"
		port = DEFAULT_PORT
	} else {
		port = appEnv.Port
		if postgres, err := appEnv.Services.WithLabel("postgresql-9.1"); err == nil {
			connString = fmt.Sprintf("%s?sslmode=disable", postgres[0].Credentials["uri"])
		}
	}

	db = NewDBConnection(connString)

	mux := goji.NewMux()
	mux.HandleFuncC(pat.Get("/"), handleIndex)
	mux.HandleFuncC(pat.Get("/api/comments"), getComments)
	mux.HandleFuncC(pat.Post("/api/comments"), postComments)
	staticPath, _ := filepath.Abs("./static")
	mux.HandleFunc(pat.Get("/static/*"), http.StripPrefix("/static", http.FileServer(http.Dir(staticPath))).ServeHTTP)

	log.Printf("Starting listening port %d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
