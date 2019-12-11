package app

import (
	"github.com/r2dev/job-search/models"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUserPost handle post /register
func (app *App) RegisterUserPost(w http.ResponseWriter, r *http.Request) {
	session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) == 0 {
		session.AddFlash("Please enter username")
		session.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	userExisted, err := app.DB.CheckUserExistByUsername(username)
	if err != nil {
		app.L.WithError(err).Warn("CheckUserExistByUsername failed")
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	if userExisted {
		session.AddFlash("User existed")
		session.Save(r, w)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	err = app.DB.CreateUserWithUsernameAndPassword(username, password)
	if err != nil {
		session.AddFlash("Something is wrong")
		session.Save(r, w)
		app.L.WithError(err)
		http.Redirect(w, r, "/register", http.StatusFound)
		return
	}
	session.AddFlash("Account created")
	session.Save(r, w)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	return
}

// LoginUserPost handle post /login
func (app *App) LoginUserPost(w http.ResponseWriter, r *http.Request) {
	session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	var user models.User
	err := app.DB.GetUserByUsername(&user, username)
	if err != nil {
		session.AddFlash("username or password is not correct")
		session.Save(r, w)
		if err == models.NoFoundUser {
			app.L.WithError(err).Info("dont get user")

		}
		http.Redirect(w, r, "login", http.StatusFound)
		return
	}
	encryptedPassword := user.Password
	err = bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(password))
	if err != nil {
		session.AddFlash("username or password is not correct")
		session.Save(r, w)
		if err == bcrypt.ErrMismatchedHashAndPassword {
			app.L.WithError(err).Info("mismatch hash")
			http.Redirect(w, r, "login", http.StatusFound)
			return
		}
		app.L.WithError(err).Info("compare hash unknwon error")
		http.Redirect(w, r, "login", http.StatusFound)
		return

	}

	session.Values["n_0"] = user.UserID.Hex()
	session.Options.HttpOnly = true
	session.Options.MaxAge = 0
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}

// LogoutUserPost handle post /logout
func (app *App) LogoutUserPost(w http.ResponseWriter, r *http.Request) {
	session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
	session.Options.MaxAge = -1
	session.AddFlash("Sign out success")
	session.Save(r, w)
	err := session.Save(r, w)
	if err != nil {
		app.L.Warning("session delete failed")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
	return
}
