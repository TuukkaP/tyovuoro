package models

type Place struct {
	Id      int64  `db:"id" json:"id"`
	Name    string `db:"name" json:"name"`
	Info    string `db:"info" json:"info"`
	Address string `db:"address" json:"address"`
}

func (p Place) GetStructMap() *map[string]string {
	return &map[string]string{
		"name":    p.Name,
		"info":    p.Info,
		"address": p.Address}
}
