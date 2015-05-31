package models

type User struct {
	UserId    int64  `db:"id" json:"id,omitempty"`
	Username  string `db:"username" json:"username,omitempty"`
	Password  string `db:"password" json:"password,omitempty"`
	Role      string `db:"role" json:"role,omitempty"`
	Address   string `db:"address" json:"address,omitempty"`
	Firstname string `db:"firstname" json:"firstname,omitempty"`
	Lastname  string `db:"lastname" json:"lastname,omitempty"`
	Email     string `db:"email" json:"email,omitempty"`
}

func (u User) GetStructMap() *map[string]interface{} {
	return &map[string]interface{}{
		"username":  u.Username,
		"password":  u.Password,
		"role":      u.Role,
		"address":   u.Address,
		"firstname": u.Firstname,
		"lastname":  u.Lastname,
		"email":     u.Email}
}

func (u User) GetId() int64 {
	return u.UserId
}
