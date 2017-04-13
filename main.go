package main

import (
	"net/http"
	"time"

	"fmt"

	"log"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))
var router = mux.NewRouter()

func main() {
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("f:/Music"))))

	router.HandleFunc("/", indexPageHandler)
	router.HandleFunc("/internal", internalPageHandler)
	router.Handle("/music", http.FileServer(http.Dir("f:/Music")))

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")

	//http.Handle("/", router)
	//log.Fatal(http.ListenAndServe(":8020", nil))

	s := &http.Server{
		Addr:           ":8022",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

func serveMusicHandler(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("f:/Music"))
}

const indexPage = `<h1>login</h1>
	<form method="post" action="/login">
		<label for="name">User name</label>
		<input type="text" id="name" name="name">
		<label for="password">Password</label>
		<input type="password" id="password" name="password">
		<button type="submit">Login</button>
	</form>`

func indexPageHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Fprintf(response, indexPage)
}

const internalPage = `
	<h1>Internal</h1>
	<hr>
		<small>User: %s</small>
	<form method="post" action="/logout">
		<button type="submit">Logout</button>
	</form>
`

func internalPageHandler(response http.ResponseWriter, request *http.Request) {
	userName := getUserName(request)
	if userName != "" {
		fmt.Fprintf(response, internalPage, userName)
	} else {
		http.Redirect(response, request, "/", 302)
	}
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
		// .. check credentials ..
		if name == "anon" && pass == "parola23" {
			setSession(name, response)
			redirectTarget = "/internal"
		}
	}
	http.Redirect(response, request, redirectTarget, 302)
}

func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(response, cookie)
	}
}

func getUserName(request *http.Request) (userName string) {
	if cookie, err := request.Cookie("session"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["name"]
		}
	}
	return userName
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}
