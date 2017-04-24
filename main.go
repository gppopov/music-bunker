package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abbot/go-http-auth"
)

func Secret(user, realm string) string {
	users := map[string]string{
		"john": "$1$dlPL2MqE$oQmn16q49SqdmhenQuNgs1", //hello
	}

	if a, ok := users[user]; ok {
		return a
	}
	return ""
}

func doRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "<h1>static file server</h1><p><a href='./static'>Music</p>")
}

func handleFileServer(dir, prefix string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))
	realHandler := http.StripPrefix(prefix, fs).ServeHTTP
	return func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.URL)
		realHandler(w, req)
	}
}

func main() {
	authenticator := auth.NewBasicAuthenticator("localhost", Secret)
	http.HandleFunc("/static/", auth.JustCheck(authenticator, handleFileServer("f:/Music", "/static/")))
	http.HandleFunc("/", auth.JustCheck(authenticator, handleFileServer("./", "/")))

	log.Println(`Listening... http://localhost:8022
 		folder is ./static
 		authentication in map users`)

	s := &http.Server{
		Addr:           ":8022",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
