package main

import (
	"encoding/json"
	"flag"
	. "github.com/TuukkaP/tyovuoro/controllers"
	"github.com/TuukkaP/tyovuoro/datastore"
	"github.com/gorilla/context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var users *UserController
var ctrl map[string]Controller
var ds = datastore.NewDatastore()
var loginCtrl = &LoginController{ds}
var sessionCtrl = &SessionController{}

func main() {
	port := flag.Int("port", 4000, "Port number")
	flag.Parse()

	ctrl = map[string]Controller{
		"users":  &UserController{ds},
		"places": &PlaceController{ds},
		"orders": &OrderController{ds},
	}
	fs := http.FileServer(http.Dir("public/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/api/", ApiResolver)
	http.HandleFunc("/", Index)

	log.Println("Server is starting in port: ", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), context.ClearHandler(http.DefaultServeMux)))
}

func ApiResolver(w http.ResponseWriter, r *http.Request) {
	username := sessionCtrl.GetUserName(w, r)
	if v := sessionCtrl.ValidSession(w, r); v == "false" {
		log.Println("Server: ApiResolver: " + username)
		http.Redirect(w, r, "/login", 302)
		return
	}
	start := time.Now()
	var response *Response
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
	w.Header().Set("content-type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	url := strings.Split(r.URL.Path, "/")
	c, ok := ctrl[url[2]]
	if ok != true {
		http.Error(w, "Resource does not exist", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		if len(url) == 4 {
			response = c.Get(w, r)
		} else {
			response = c.GetAll(w, r)
		}
	case "PUT":
		response = c.Update(w, r)
	case "POST":
		response = c.Create(w, r)
	case "DELETE":
		response = c.Delete(w, r)
	default:
		http.Error(w, "Wrong http method", http.StatusMethodNotAllowed)
		return
	}

	// Handle reponse errors
	if response.Err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		log.Println(response.Err)
		return
	}

	bytes, e := json.Marshal(response.Data)
	if e != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
	log.Println(username, r.Method, r.RequestURI, time.Since(start))
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("---------Server: Login START--------------")
	switch r.Method {
	case "POST":
		response := loginCtrl.Login(w, r, sessionCtrl)
		if response.Err == nil && response.Data != nil {
			log.Println("Server: Login redirect")
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			log.Println("Server: Login serve file")
			http.ServeFile(w, r, "public/login.html")
		}
	default:
		http.ServeFile(w, r, "public/login.html")
	}
	log.Println("---------Server: Login END --------------")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	log.Println("---------Server: Logout START--------------")
	loginCtrl.Logout(w, r, sessionCtrl)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
	log.Println("---------Server: Logout END--------------")
}

func Index(w http.ResponseWriter, r *http.Request) {
	log.Println("---------Server: Index START--------------")
	if v := sessionCtrl.ValidSession(w, r); v == "true" {
		log.Println("Server: Index: " + sessionCtrl.GetUserName(w, r) + " Valid: " + v)
		http.ServeFile(w, r, "public/index.html")
	} else {
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	log.Println("---------Server: Index END--------------")
}
