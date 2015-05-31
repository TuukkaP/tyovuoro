package datastore

import (
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	. "github.com/TuukkaP/tyovuoro/models"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	query := fmt.Sprintf("SELECT * FROM %v %v", table, string(table[0]))
	var stmt *sqlx.Stmt
	var err error
	if table == "orders" {
		return d.parseOrders(model, id)
	}
	if id != nil {
		query += fmt.Sprintf(" where %v.id = $1", string(table[0]))
		stmt, err = d.Db.Preparex(query)
		stmt.Select(model, id)
	} else {
		stmt, err = d.Db.Preparex(query)
		stmt.Select(model)
	}
	return err
}

func (d Datastore) Update(table string, model Model, url_id int64, denied_fields ...interface{}) error {
	query := "UPDATE " + table + " SET"
	id := model.GetId()
	/* var id int64
	 val := reflect.ValueOf(model).Elem()
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
	}*/
	for k, v := range *model.GetStructMap() {
		if v != "" {
			query += fmt.Sprintf(" %v = '%v',", k, v)
		}
	}

	// Sanity check that the POST and url id match
	if id == url_id {
		query = query[:len(query)-1] + " WHERE id = " + strconv.FormatInt(id, 10)
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
	}
	var names string
	for k := range *model.GetStructMap() {
		names += fmt.Sprintf(":%v,", k)
	}
	names = names[:len(names)-1]
	query := fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", table, strings.Replace(names, ":", "", -1), names)
	log.Println(query)
	log.Println(model)
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

// Orders need to be joined to extract usefull info, therefore they have to be parsed by hand
func (d Datastore) parseOrders(model interface{}, id interface{}) error {
	var rows *sqlx.Rows
	var e error
	query := "select o.id, o.user_id, o.place_id, o.start, o.end_time, p.name as place_name, u.username, u.firstname, u.lastname from orders o left join users u on u.id = o.user_id left join places p on p.id = o.place_id"
	/*	query := "select o.id, o.user_id, o.place_id, o.date + o.order_start as start, o.date + o.order_end as end, p.name as place_name, u.username, u.firstname, u.lastname from orders o left join users u on u.id = o.user_id left join places p on p.id = o.place_id"
	 */if id != nil {
		query += " where o.id = $1"
		rows, e = d.Db.Queryx(query, id)
	} else {
		rows, e = d.Db.Queryx(query)
	}
	/*query = "SELECT * FROM orders o, users u, places p where u.id = o.user_id and p.id = o.place_id"*/
	if e != nil {
		log.Println(e)
	}
	for rows.Next() {
		var id int64
		var place_id, user_id sql.NullInt64
		var username, firstname, lastname, place_name sql.NullString
		var start, end pq.NullTime
		if err := rows.Scan(&id, &user_id, &place_id, &start, &end, &place_name, &username, &firstname, &lastname); err != nil {
			log.Println(err)
		}
		order := Order{Id: id, UserId: user_id.Int64, PlaceId: place_id.Int64, Start: start.Time, End: end.Time, PlaceName: place_name.String, Username: username.String, Firstname: firstname.String, Lastname: lastname.String}
		if user_id.Valid == false || (lastname.String == "" && firstname.String == "") {
			order.Title = fmt.Sprintf("%v", place_name.String)
		} else {
			order.Title = fmt.Sprintf("%v: %v, %v", place_name.String, lastname.String, firstname.String)
		}
		*model.(*[]Order) = append(*model.(*[]Order), order)
	}
	e = rows.Err()
	return e
}

func (d Datastore) Login(name string, pass string, user *User) error {
	hash := sha512.New()
	hash.Write([]byte(pass))
	password := fmt.Sprintf("%v", base64.URLEncoding.EncodeToString(hash.Sum(nil)))
	stmt, err := d.Db.Preparex("select * from users where username = $1 and password = $2")
	stmt.Get(user, name, password)
	return err
}
