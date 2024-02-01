package main

import (
	"errors"
	"fmt"
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

var users = make(map[string]User)

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
	c, err := r.Cookie("session")
	if errors.Is(err, http.ErrNoCookie) {
		tpl.ExecuteTemplate(w, "index.gohtml", nil)
		return
	}

	u := users[c.Value]
	fmt.Println(users)
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

		users[sid] = nu
		fmt.Println(users)
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func bar(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if errors.Is(err, http.ErrNoCookie) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	u := users[c.Value]
	fmt.Println(users)
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}
