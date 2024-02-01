package main

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
)

func getUser(w http.ResponseWriter, r *http.Request) User {
	c, err := r.Cookie("session")
	if errors.Is(err, http.ErrNoCookie) {
		sid := uuid.NewString()
		http.SetCookie(w, &http.Cookie{Name: "session", Value: sid, HttpOnly: true})
	}

	var u User

	if sid, ok := dbSessions[c.Value]; ok {
		u = dbUsers[sid]
	}
	return u
}

func alreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if errors.Is(err, http.ErrNoCookie) {
		return false
	}

	_, ok := dbSessions[c.Value]
	return ok
}
