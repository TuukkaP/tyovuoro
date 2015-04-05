package main

import (
	"encoding/json"
	"flag"
	. "github.com/TuukkaP/tyovuoro/controllers"
	"github.com/TuukkaP/tyovuoro/datastore"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var users *UserController
var ctrl map[string]Controller

func main() {
	port := flag.Int("port", 4000, "Port number")
	flag.Parse()

	/*	log.Println("Connecting to postgres")
		db, err := sqlx.Open("postgres", "user=tuukka password=tuukka port=5433 dbname=peuranie sslmode=disable")
		if err != nil {
			log.Println(err)
		}
		log.Println(db.Ping())

		ds := &datastore.Datastore{db}*/
	ds := datastore.NewDatastore()
	/*users = &UserController{ds}*/
	ctrl = map[string]Controller{
		"users":  &UserController{ds},
		"places": &PlaceController{ds},
	}

	http.HandleFunc("/", http.HandlerFunc(ServeIndex))
	/*	http.HandleFunc("/users", ApiHandler(users.GetAll, "GET"))
		http.HandleFunc("/users/", ApiHandler(users.Get, "GET"))
		http.HandleFunc("/usersa/", ApiHandler(users.Update, "PUT"))*/
	http.HandleFunc("/api/", ApiResolver)

	log.Println("Server is starting in port: ", *port)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(*port), nil))
}

func ApiResolver(w http.ResponseWriter, r *http.Request) {
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
		switch {
		case len(url) == 3:
			response = c.GetAll(w, r)
		case len(url) == 4:
			response = c.Get(w, r)
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
	log.Println(r.Method, r.RequestURI, time.Since(start))
}

func ApiHandler(fn func(http.ResponseWriter, *http.Request) *Response, method string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check for allowed HTTP method
		if r.Method != method {
			http.Error(w, "Wrong http method", http.StatusMethodNotAllowed)
			return
		}

		start := time.Now()
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("content-type", "application/json; charset=utf-8")

		response := fn(w, r)

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

		log.Println(r.Method, r.RequestURI, time.Since(start))
		w.Write(bytes)
	}
}

func ServeIndex(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method + ": " + r.RequestURI)
	http.ServeFile(w, r, "views/index.html")
}
