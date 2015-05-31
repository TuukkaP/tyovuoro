package controllers

import (
	"errors"
	"github.com/TuukkaP/tyovuoro/datastore"
	. "github.com/TuukkaP/tyovuoro/models"
	"log"
	"net/http"
)

type LoginController struct {
	Datastore *datastore.Datastore
}

func (lc LoginController) Login(w http.ResponseWriter, r *http.Request, sessionCtrl *SessionController) *Response {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	if name != "" && pass != "" {
		log.Println(r.Method + ": " + r.RequestURI + "[" + name + "]")
		user := User{}
		err := lc.Datastore.Login(name, pass, &user)
		if err == nil && user.Username != "" && user.UserId != 0 {
			log.Printf("%v is logging in!", user.Username)
			sessionCtrl.SetSession(w, r, user)
		} else {
			sessionCtrl.SetFlash(w, r, "Authentication failed!")
			log.Println(err)
		}
		return &Response{user, nil}
	} else {
		return &Response{nil, errors.New("Username or password was empty")}
	}
}
