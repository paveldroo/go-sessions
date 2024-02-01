package main

import (
	"github.com/google/uuid"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	UserName string
	First    string
	Last     string
}

var tpl *template.Template

var dbUsers = make(map[string]User)
var dbSessions = make(map[string]string)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/bar", bar)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		tpl.ExecuteTemplate(w, "index.gohtml", nil)
		return
	}
	u := getUser(w, r)
	tpl.ExecuteTemplate(w, "index.gohtml", u)
}

func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		sid := uuid.NewString()
		http.SetCookie(w, &http.Cookie{Name: "session", Value: sid, HttpOnly: true})
		nu := User{
			UserName: r.FormValue("username"),
			First:    r.FormValue("firstname"),
			Last:     r.FormValue("lastname"),
		}
		dbSessions[sid] = nu.UserName
		dbUsers[nu.UserName] = nu
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func bar(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	u := getUser(w, r)
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}
