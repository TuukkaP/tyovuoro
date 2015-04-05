package models

import (
	"time"
)

type Order struct {
	Id        int64 `db:"id" json:"id"`
	Place     Place
	User      User
	Date      time.Time `db:"date" json:"date"`
	StartDate time.Time `db:"order_start" json:"order_start"`
	EndDate   time.Time `db:"order_end" json:"order_end"`
}

func (o Order) GetStructMap() *map[string]string {
	return &map[string]string{
		"date":       o.Date.String(),
		"start_date": o.StartDate.String(),
		"end_date":   o.EndDate.String()}
}
