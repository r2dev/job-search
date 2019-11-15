package app

import (
	"context"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/tj/go/http/response"
)

type ctxKeyUserID struct{}
type ctxKeySession struct{}

func (app *App) WithUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value(ctxKeySession{}).(*sessions.Session)
		if !ok {
			response.InternalServerError(w)
		}
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

func (app *App) WithSession(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := app.S.Get(r, "r_u_n_a_w_a_y")
		ctx := context.WithValue(r.Context(), ctxKeySession{}, session)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
