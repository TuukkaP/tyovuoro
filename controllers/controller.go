package controllers

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Data interface{}
	Err  error
}

// CRUD interface
type Controller interface {
	GetAll(w http.ResponseWriter, r *http.Request) *Response
	Get(w http.ResponseWriter, r *http.Request) *Response
	Update(w http.ResponseWriter, r *http.Request) *Response
	Create(w http.ResponseWriter, r *http.Request) *Response
	Delete(w http.ResponseWriter, r *http.Request) *Response
}

func parseId(url string) (int64, error) {
	id, err := strconv.ParseInt(strings.Split(url, "/")[3], 10, 32)
	if err != nil {
		log.Println("Url fields", url)
		log.Println("strconv error", err)
		return 0, err
	}
	return id, nil
}

func logError(obj interface{}, err error) interface{} {
	log.Println("Error occurred!", err)
	return obj
}
