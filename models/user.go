package models

type User struct {
	Id        int64  `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	Password  string `db:"password" json:"password"`
	Role      string `db:"role" json:"role"`
	Address   string `db:"address" json:"address"`
	Firstname string `db:"firstname" json:"firstname"`
	Lastname  string `db:"lastname" json:"lastname"`
	Email     string `db:"email" json:"email"`
}

func (u User) GetStructMap() *map[string]string {
	return &map[string]string{
		"username":  u.Username,
		"password":  u.Password,
		"role":      u.Role,
		"address":   u.Address,
		"firstname": u.Firstname,
		"lastname":  u.Lastname,
		"email":     u.Email}
}
