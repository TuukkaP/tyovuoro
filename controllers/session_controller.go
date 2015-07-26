package controllers

import (
	"fmt"
	. "github.com/TuukkaP/tyovuoro/models"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"time"
)

type SessionController struct{}

var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(64)))

func (sc SessionController) SetSession(w http.ResponseWriter, r *http.Request, user User) {
	session := sc.GetSession(w, r)
	session.Values["username"] = user.Username
	session.Values["user_id"] = user.UserId
	session.Values["timestamp"] = time.Now().String()
	session.Values["valid"] = "true"
	uc := http.Cookie{Name: "username", Value: user.Username}
	idc := http.Cookie{Name: "user_id", Value: fmt.Sprintf("%v", user.UserId)}
	http.SetCookie(w, &uc)
	http.SetCookie(w, &idc)
	sc.SaveSession(w, r, session)
	log.Println("Cookie was set for " + user.Username)
}

func (sc SessionController) SetFlash(w http.ResponseWriter, r *http.Request, msg string) {
	session := sc.GetSession(w, r)
	session.AddFlash(msg)
	sc.SaveSession(w, r, session)
}

func (sc SessionController) ClearSession(w http.ResponseWriter, r *http.Request) {
	session := sc.GetSession(w, r)
	session.Values["username"] = nil
	session.Values["user_id"] = nil
	session.Values["timestamp"] = nil
	session.Values["valid"] = "false"
	session.AddFlash("Logout!")
	session.Options = &sessions.Options{
		MaxAge: -1,
	}
	sc.SaveSession(w, r, session)
}

func (sc SessionController) GetUserName(w http.ResponseWriter, r *http.Request) string {
	session := sc.GetSession(w, r)
	userName, ok := session.Values["username"].(string)
	if ok == false {
		userName = ""
	}
	log.Println("SessionController: GetUserName: " + userName)
	return userName
}

func (sc SessionController) GetUserId(w http.ResponseWriter, r *http.Request) int64 {
	session := sc.GetSession(w, r)
	id, ok := session.Values["user_id"].(int64)
	if ok == false {
		id = -1
	}
	log.Println("SessionController: GetUserId: " + fmt.Sprintf("%v", id))
	return id
}

func (sc SessionController) GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, err := store.Get(r, "tyovuoro-sessio")
	if err != nil {
		log.Println("ERROR: SessionController: GetSession: " + err.Error())
	}
	return session
}

func (sc SessionController) SaveSession(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	err := s.Save(r, w)
	if err != nil {
		log.Println("ERROR: SessionController: SaveSession: " + err.Error())
	}
}

func (sc SessionController) ValidSession(w http.ResponseWriter, r *http.Request) string {
	session := sc.GetSession(w, r)
	valid, ok := session.Values["valid"].(string)
	if ok == false {
		valid = "false"
	}
	log.Println("SessionController: ValidSession: " + valid)
	return valid
}
