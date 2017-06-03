package controllers

import (
	"html/template"
	"net/http"

	"github.com/spf13/viper"

	log "github.com/cihub/seelog"
)

const (
	LOGIN_TEMPLATE_PATH = "views/login.html"
	HOME_TEMPLATE_PATH  = "views/home.html"
)

func Home(w http.ResponseWriter, r *http.Request) {
	if !hasSession(w, r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	renderHome(w)
}

func renderHome(w http.ResponseWriter) {
	tpl := template.Must(template.ParseFiles(HOME_TEMPLATE_PATH))
	tpl.Execute(w, nil)
}

func renderLogin(w http.ResponseWriter) {
	tpl := template.Must(template.ParseFiles(LOGIN_TEMPLATE_PATH))
	tpl.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Info("login")
	if r.Method == http.MethodGet {
		renderLogin(w)
		return
	}
	user := r.FormValue("username")
	password := r.FormValue("userpassword")

	if user == viper.GetString("login.user") && password == viper.GetString("login.password") {
		setSession(user, w, r)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	renderLogin(w)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	closeSession(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
