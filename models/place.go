package models

type Place struct {
	PlaceId int64  `db:"id" json:"id,omitempty"`
	Name    string `db:"name" json:"name,omitempty"`
	Info    string `db:"info" json:"info,omitempty"`
	Address string `db:"address" json:"address,omitempty"`
}

func (p Place) GetStructMap() *map[string]interface{} {
	return &map[string]interface{}{
		"name":    p.Name,
		"info":    p.Info,
		"address": p.Address}
}

func (p Place) GetId() int64 {
	return p.PlaceId
}
