package controllers

import (
	. "github.com/TuukkaP/tyovuoro/models"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"log"
	"net/http"
	"time"
)

type SessionController struct{}

/*var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))*/
var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(64)))

func (sc SessionController) SetSession(w http.ResponseWriter, r *http.Request, user User) {
	session := sc.GetSession(w, r)
	session.Values["username"] = user.Username
	session.Values["user_id"] = user.UserId
	session.Values["timestamp"] = time.Now().String()
	session.AddFlash("Welcome!", "flash")
	sc.SaveSession(w, r, session)
	log.Println("Cookie was set for " + user.Username)
	/*	value := map[string]string{
		"username": user.Username,
		"id":       strconv.FormatInt(user.UserId, 10),
	}*/
	/*	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		cookie := &http.Cookie{
			Name:  "Tyovuoro",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	}*/
}

func (sc SessionController) SetFlash(w http.ResponseWriter, r *http.Request, msg string) {
	session := sc.GetSession(w, r)
	session.AddFlash(msg)
	sc.SaveSession(w, r, session)
}

func (sc SessionController) ClearSession(w http.ResponseWriter, r *http.Request) {
	/*	cookie := &http.Cookie{
			Name:   "Tyovuoro",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)*/
	session := sc.GetSession(w, r)
	session.Values["username"] = nil
	session.Values["user_id"] = nil
	session.Values["timestamp"] = nil
	session.AddFlash("Logout!")
	sc.SaveSession(w, r, session)
}

func (sc SessionController) GetUserName(r *http.Request) string {
	/*	if cookie, err := request.Cookie("Tyovuoro"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			userName = cookieValue["username"]
		}
	}*/
	session := sc.GetSession(w, r)
	userName, ok := session.Values["username"].(string)
	if ok == false {
		userName = ""
		log.Println(err)
	}
	return userName
}

func (sc SessionController) GetUserId(r *http.Request) int64 {
	/*	if cookie, err := request.Cookie("Tyovuoro"); err == nil {
		cookieValue := make(map[string]string)
		if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
			id, _ = strconv.ParseInt(cookieValue["id"], 10, 64)
		}
	}*/
	session := sc.GetSession(w, r)
	id, ok := session.Values["user_id"].(int64)
	if ok == false {
		id = -1
		log.Println(err)
	}
	return id
}

func (sc SessionController) GetSession(w http.ResponseWriter, r *http.Request) *sessions.Session {
	session, err := store.Get(r, "tyovuoro-sessio")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return session
}

func (sc SessionController) SaveSession(w http.ResponseWriter, r *http.Request, s *sessions.Session) {
	err := session.Save(r, w)
	if err != nil {
		log.Println(err)
	}
}
