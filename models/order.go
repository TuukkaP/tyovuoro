package models

import (
	"fmt"
	"time"
)

type Order struct {
	Id        int64     `db:"id" json:"id"`
	PlaceId   int64     `db:"place_id" json:"place_id"`
	UserId    int64     `db:"user_id" json:"user_id"`
	Title     string    `json:"title"`
	Start     time.Time `db:"start" json:"start"`
	End       time.Time `db:"end_time" json:"end_time"`
	Username  string    `json:"username"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	PlaceName string    `json:"place_name"`
}

func (o Order) GetStructMap() *map[string]interface{} {
	order := make(map[string]interface{})

	if o.PlaceId != 0 {
		order["place_id"] = fmt.Sprint(o.PlaceId)
	}

	if o.UserId != 0 {
		order["user_id"] = fmt.Sprint(o.UserId)
	}

	if year := fmt.Sprint(o.Start.Year()); year != "1" {
		order["start"] = o.Start.String()[:len(o.Start.String())-10]
	}

	if year := fmt.Sprint(o.End.Year()); year != "1" {
		order["end_time"] = o.End.String()[:len(o.End.String())-10]
	}
	return &order
}

func (o Order) GetId() int64 {
	return o.Id
}
