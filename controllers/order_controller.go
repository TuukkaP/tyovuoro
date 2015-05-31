package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/TuukkaP/tyovuoro/datastore"
	"github.com/TuukkaP/tyovuoro/models"
	"log"
	"net/http"
)

type OrderController struct {
	Datastore *datastore.Datastore
}

func (oc OrderController) GetAll(w http.ResponseWriter, r *http.Request) *Response {
	orders := []models.Order{}
	err := oc.Datastore.Get("orders", &orders, nil)
	log.Println(orders)
	return &Response{orders, err}
}

func (oc OrderController) Get(w http.ResponseWriter, r *http.Request) *Response {
	orders := []models.Order{}
	id, err := parseId(r.URL.Path)
	if err != nil {
		return &Response{nil, err}
	}
	err = oc.Datastore.Get("orders", &orders, id)
	return &Response{orders, err}
}

func (oc OrderController) Create(w http.ResponseWriter, r *http.Request) *Response {
	var user models.Order
	err := json.NewDecoder(r.Body).Decode(&user)
	fmt.Printf("%+v\n", r)
	if err != nil {
		log.Println("Error parsing user from JSON", err)
	} else {
		err = oc.Datastore.Create("orders", &user)
	}
	fmt.Printf("%+v\n", user)
	log.Println(err)
	return &Response{nil, err}
}

func (oc OrderController) Update(w http.ResponseWriter, r *http.Request) *Response {
	var user models.Order
	err := json.NewDecoder(r.Body).Decode(&user)
	fmt.Printf("%+v\n", r)
	if err != nil {
		log.Println("Error parsing user from JSON", err)
	}
	fmt.Printf("%+v\n", user)
	id, err := parseId(r.URL.Path)
	log.Println(id)
	if err == nil {
		err = oc.Datastore.Update("orders", &user, id, "id", "password")
	}
	log.Println(err)
	return &Response{nil, err}
}

func (oc OrderController) Delete(w http.ResponseWriter, r *http.Request) *Response {
	id, err := parseId(r.URL.Path)
	if err != nil {
		return &Response{nil, err}
	}
	err = oc.Datastore.Delete("orders", id)
	return &Response{nil, err}
}
