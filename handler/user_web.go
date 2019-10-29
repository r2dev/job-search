package handler

import (
	"hirine/models"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	// http.Redirect(w, r, "/register", 302)
	_, err := models.GetUserByUsername(username)
	if err != nil && err != models.NoFoundUser {
		log.WithError(err)
		http.Redirect(w, r, "/register", 302)
		return
	}
	_, err = models.CreateUserWithUsernameAndPassword(username, password)
	if err != nil {
		log.WithError(err)
		http.Redirect(w, r, "/register", 302)
		return
	}
	// http.Redirect(w, r, "")
	// response.OK(w, map[string]string{
	// 	"id": id,
	// })
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// username := r.FormValue("username")
	// password := r.FormValue("password")

}
