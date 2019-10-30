package app

import (
	"hirine/models"
	"net/http"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

func (app *App) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	_, err := app.DB.GetUserByUsername(username)
	if err != nil && err != models.NoFoundUser {
		log.WithError(err)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	_, err = app.DB.CreateUserWithUsernameAndPassword(username, password)
	if err != nil {
		log.WithError(err)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *App) LoginHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := app.DB.GetUserByUsername(username)
	if err != nil {
		if err == models.NoFoundUser {
			log.WithError(err).Info("dont get user")
		}
		http.Redirect(w, r, "login", http.StatusFound)
		return
	}
	encryptedPassword := user.Password
	err = bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			log.WithError(err).Info("mismatch hash")
			http.Redirect(w, r, "login", http.StatusFound)
			return
		} else {
			log.WithError(err).Info("compare hash unknwon error")
			http.Redirect(w, r, "login", http.StatusFound)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return

}
