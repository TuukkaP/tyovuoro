package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/TuukkaP/tyovuoro/datastore"
	. "github.com/TuukkaP/tyovuoro/models"
	"log"
	"net/http"
)

type PlaceController struct {
	Datastore *datastore.Datastore
}

func (pc PlaceController) GetAll(w http.ResponseWriter, r *http.Request) *Response {
	places := []Place{}
	err := pc.Datastore.Get("places", &places, nil)
	return &Response{places, err}
}

func (pc PlaceController) Get(w http.ResponseWriter, r *http.Request) *Response {
	places := []Place{}
	id, err := parseId(r.URL.Path)
	if err != nil {
		return &Response{nil, err}
	}
	err = pc.Datastore.Get("places", &places, id)
	return &Response{places, err}
}

func (pc PlaceController) Create(w http.ResponseWriter, r *http.Request) *Response {
	var place Place
	err := json.NewDecoder(r.Body).Decode(&place)
	fmt.Printf("%+v\n", r)
	if err != nil {
		log.Println("Error parsing place from JSON", err)
	} else {
		err = pc.Datastore.Create("places", &place)
	}
	fmt.Printf("%+v\n", place)
	log.Println(err)
	return &Response{nil, err}
}

func (pc PlaceController) Update(w http.ResponseWriter, r *http.Request) *Response {
	var place Place
	err := json.NewDecoder(r.Body).Decode(&place)
	fmt.Printf("%+v\n", r)
	if err != nil {
		log.Println("Error parsing place from JSON", err)
	}
	fmt.Printf("%+v\n", place)
	id, err := parseId(r.URL.Path)
	log.Println(id)
	if err == nil {
		err = pc.Datastore.Update("places", &place, id, "id", "password")
	}
	log.Println(err)
	return &Response{nil, err}
}

func (pc PlaceController) Delete(w http.ResponseWriter, r *http.Request) *Response {
	id, err := parseId(r.URL.Path)
	if err != nil {
		return &Response{nil, err}
	}
	err = pc.Datastore.Delete("places", id)
	return &Response{nil, err}
}
