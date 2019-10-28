package handler

import (
	"net/http"
	"text/template"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	_, uok := r.Form["username"]
	_, pok := r.Form["password"]
	if !uok || !pok {
		var indexTemp = template.Must(
			template.ParseFiles("./templates/layout/base.html", "./templates/register.html"))
		indexTemp.Execute(w, nil)
		return
	}
	// username := request.Username
	// password := request.Password
	// _, err = models.GetUserByUsername(username)
	// if err != nil {
	// 	if err != models.NoFoundUser {
	// 		log.WithError(err)
	// 		response.Error(w, http.StatusInternalServerError)
	// 		return
	// 	}
	// }
	// id, err := models.CreateUserWithUsernameAndPassword(username, password)
	// if err != nil {
	// 	log.WithError(err)
	// 	response.Error(w, http.StatusBadGateway)
	// 	return
	// }
	// response.OK(w, map[string]string{
	// 	"id": id,
	// })

}
