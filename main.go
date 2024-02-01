package main

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"log"
	"net/http"
)

type User struct {
	UserName string
	Password []byte
	First    string
	Last     string
	Role     string
}

var tpl *template.Template

var dbUsers = make(map[string]User)
var dbSessions = make(map[string]string)

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	p, err := bcrypt.GenerateFromPassword([]byte("123"), bcrypt.MinCost)
	if err != nil {
		log.Fatalln(err)
	}
	nu := User{
		UserName: "123@123.ru",
		Password: p,
		First:    "Yo",
		Last:     "YOYO",
		Role:     "user",
	}
	dbUsers[nu.UserName] = nu
}

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
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

func login(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		p := r.FormValue("password")
		if u, ok := dbUsers[un]; ok {
			err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
			if err != nil {
				http.Error(w, "Incorrect login or password", http.StatusForbidden)
				return
			}
			sid := uuid.NewString()
			http.SetCookie(w, &http.Cookie{Name: "session", Value: sid, HttpOnly: true})
			dbSessions[sid] = u.UserName
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		http.Error(w, "Incorrect login or password", http.StatusForbidden)
		return
	}
	tpl.ExecuteTemplate(w, "login.gohtml", nil)
}

func logout(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	c, _ := r.Cookie("session")
	delete(dbSessions, c.Value)
	http.SetCookie(w, &http.Cookie{Name: "session", Value: "", MaxAge: -1})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func signup(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		f := r.FormValue("firstname")
		l := r.FormValue("lastname")
		p := r.FormValue("username")
		role := r.FormValue("role")

		if _, ok := dbUsers[un]; ok {
			http.Error(w, "Username already taken", http.StatusForbidden)
			return
		}

		sid := uuid.NewString()
		http.SetCookie(w, &http.Cookie{Name: "session", Value: sid, HttpOnly: true})
		ep, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "There was an error with user data", http.StatusInternalServerError)
		}
		nu := User{
			UserName: un,
			Password: ep,
			First:    f,
			Last:     l,
			Role:     role,
		}
		fmt.Println(string(nu.Password))

		dbSessions[sid] = nu.UserName
		dbUsers[nu.UserName] = nu
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	tpl.ExecuteTemplate(w, "signup.gohtml", nil)
}

func bar(w http.ResponseWriter, r *http.Request) {
	if !alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if !hasAccess(r) {
		http.Redirect(w, r, "/", http.StatusForbidden)
		return
	}

	u := getUser(w, r)
	tpl.ExecuteTemplate(w, "bar.gohtml", u)
}
