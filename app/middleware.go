package app

import (
	"context"
	"net/http"
)

type ctxKeyUserID struct{}

func (app *App) WithUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		userID, ok := session.Values["n_0"].(string)
		if !ok {
			session.AddFlash("Please login first")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		ctx := context.WithValue(r.Context(), ctxKeyUserID{}, userID)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
