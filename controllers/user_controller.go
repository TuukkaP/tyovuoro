package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/TuukkaP/tyovuoro/datastore"
	"github.com/TuukkaP/tyovuoro/models"
	"log"
	"net/http"
)

type UserController struct {
	Datastore *datastore.Datastore
}

func (uc UserController) GetAll(w http.ResponseWriter, r *http.Request) *Response {
	users := []models.User{}
	err := uc.Datastore.Get("users", &users, nil)
	return &Response{users, err}
}

func (uc UserController) Get(w http.ResponseWriter, r *http.Request) *Response {
	users := []models.User{}
	id, err := parseId(r.URL.Path)
	if err != nil {
		return &Response{nil, err}
	}
	err = uc.Datastore.Get("users", &users, id)
	return &Response{users, err}
}

func (uc UserController) Create(w http.ResponseWriter, r *http.Request) *Response {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	fmt.Printf("%+v\n", r)
	if err != nil {
		log.Println("Error parsing user from JSON", err)
	} else {
		err = uc.Datastore.Create("users", &user)
	}
	fmt.Printf("%+v\n", user)
	log.Println(err)
	return &Response{nil, err}
}

func (uc UserController) Update(w http.ResponseWriter, r *http.Request) *Response {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	fmt.Printf("%+v\n", r)
	if err != nil {
		log.Println("Error parsing user from JSON", err)
	}
	fmt.Printf("%+v\n", user)
	id, err := parseId(r.URL.Path)
	log.Println(id)
	if err == nil {
		err = uc.Datastore.Update("users", &user, id, "id", "password")
	}
	log.Println(err)
	return &Response{nil, err}
}

func (uc UserController) Delete(w http.ResponseWriter, r *http.Request) *Response {
	id, err := parseId(r.URL.Path)
	if err != nil {
		return &Response{nil, err}
	}
	err = uc.Datastore.Delete("users", id)
	return &Response{nil, err}
}
