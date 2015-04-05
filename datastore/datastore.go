package datastore

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
	. "github.com/TuukkaP/tyovuoro/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"reflect"
	"strconv"
	"strings"
)

type Datastore struct {
	Db *sqlx.DB
}

func NewDatastore() *Datastore {
	log.Println("Connecting to postgres")
	db := sqlx.MustOpen("postgres", "user=tuukka password=tuukka port=5433 dbname=peuranie sslmode=disable")
	return &Datastore{db}
}

func (d Datastore) Get(table string, model interface{}, id interface{}) error {
	query := "SELECT * FROM " + table
	var stmt *sqlx.Stmt
	var err error
	switch {
	case id == nil:
		stmt, err = d.Db.Preparex(query)
		stmt.Select(model)
	case table == "orders":
		query = "select * from orders o, users u, places p where u.id = o.user_id and p.id = o.place_id"
		stmt, err = d.Db.Preparex(query)
		stmt.Select(model)
	default:
		query += " WHERE id=$1"
		stmt, err = d.Db.Preparex(query)
		stmt.Select(model, id)
	}
	return err
}

func (d Datastore) Update(table string, model interface{}, url_id int64, denied_fields ...interface{}) error {
	query := "UPDATE " + table + " SET"
	val := reflect.ValueOf(model).Elem()
	var id int64
	var number_fields int
	for i := 0; i < val.NumField(); i++ {
		field_name := val.Type().Field(i).Tag.Get("db")
		if !Contains(denied_fields, field_name) {
			number_fields = i
		}
	}
	for i := 0; i < val.NumField(); i++ {
		field_name := val.Type().Field(i).Tag.Get("db")
		field_value := val.Field(i).Interface()

		if !Contains(denied_fields, field_name) {
			query += fmt.Sprintf(" %v = '%v'", field_name, field_value)
			log.Printf("'%+v'", field_name)
			if i < number_fields {
				query += ", "
			}
		}

		if field_name == "id" {
			id = field_value.(int64)
		}
	}

	// Sanity check that the POST and url id match
	if id == url_id {
		query += " WHERE id = " + strconv.FormatInt(id, 10)
		log.Println(query)
		_, err := d.Db.Exec(query)
		return err
	}
	return errors.New("JSON POST id and URL id did not match!")
}

func (d Datastore) Create(table string, model Model) error {
	var err error
	log.Println(reflect.TypeOf(model))
	switch model.(type) {
	case *User:
		if len(model.(*User).Password) < 8 {
			return errors.New("Password has to be atleast 8 chars long!")
		}
		hash := sha512.New()
		hash.Write([]byte(model.(*User).Password))
		model.(*User).Password = fmt.Sprintf("%v", base64.URLEncoding.EncodeToString(hash.Sum(nil)))
		/*		query := "INSERT INTO users (username, password, role, address, firstname, lastname, email) VALUES  (:username, :password, :role, :address, :firstname, :lastname, :email)"
				_, err = d.Db.NamedExec(query, map[string]interface{}{
					"username":  model.(*User).Username,
					"password":  fmt.Sprintf("%v", base64.URLEncoding.EncodeToString(hash.Sum(nil))),
					"role":      model.(*User).Role,
					"address":   model.(*User).Address,
					"firstname": model.(*User).Firstname,
					"lastname":  model.(*User).Lastname,
					"email":     model.(*User).Email})
				log.Println(query)*/
	}
	var names string
	for k := range *model.GetStructMap() {
		names += fmt.Sprintf(":%v,", k)
	}
	names = names[:len(names)-1]
	query := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", table, strings.Replace(names, ":", "", -1), names)
	log.Println(query)
	_, err = d.Db.NamedExec(query, model)
	return err
}

func (d Datastore) Delete(table string, id int64) error {
	query := fmt.Sprintf("DELETE FROM %v WHERE id = '%v'", table, id)
	stmt, err := d.Db.Preparex(query)
	_, err = stmt.Exec()
	return err
}

func Contains(list []interface{}, elem interface{}) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}
	return false
}
