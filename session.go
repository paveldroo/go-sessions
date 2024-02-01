package main

import (
	"errors"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func getUser(w http.ResponseWriter, r *http.Request) User {
	c, err := r.Cookie("session")
	if errors.Is(err, http.ErrNoCookie) {
		sid := uuid.NewString()
		http.SetCookie(w, &http.Cookie{Name: "session", Value: sid, HttpOnly: true})
	}

	var u User

	if session, ok := dbSessions[c.Value]; ok {
		u = dbUsers[session.un]
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

func hasAccess(r *http.Request) bool {
	c, err := r.Cookie("session")
	if errors.Is(err, http.ErrNoCookie) {
		return false
	}

	var u User
	if session, ok := dbSessions[c.Value]; ok {
		u = dbUsers[session.un]
	}

	return u.Role == "007"
}

func cleanSessions() {
	for sid, session := range dbSessions {
		if time.Now().Sub(session.lastActivity) > cleanDelay {
			delete(dbSessions, sid)
		}
	}
	lastCleaned = time.Now()
}
