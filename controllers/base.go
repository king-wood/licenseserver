package controllers

import (
	"licenseserver/controllers/internalerrors"
	"net/http"

	"github.com/gorilla/sessions"

	"encoding/json"
)

var (
	store = sessions.NewCookieStore([]byte("something-very-secret"))
)

const (
	SESSION_NAME = "serial-session"
)

func handleError(w http.ResponseWriter, err *internalerrors.LogicError) {
	if err.Type == internalerrors.RequestError {
		body, _ := json.Marshal(map[string]string{
			"error": err.Description,
		})
		w.Write(body)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		body, _ := json.Marshal(map[string]string{
			"error": "server internal error",
		})
		w.Write(body)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func Response(w http.ResponseWriter, code int, content interface{}) {
	body, _ := json.Marshal(content)
	w.Write(body)
	w.WriteHeader(code)
}

func hasSession(w http.ResponseWriter, r *http.Request) bool {
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return false
	}
	if session.Values["User"] != nil {
		return true
	}
	return false
}

func setSession(user string, w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	session.Values["User"] = user
	session.Save(r, w)
}

func closeSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, SESSION_NAME)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	delete(session.Values, "User")
	session.Save(r, w)
}
