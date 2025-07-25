package models

import "database/sql"

type User struct {
	Id        string         `db:"id" json:"uuid"`
	Username  string         `db:"username" json:"username"`
	Password  string         `db:"password" json:"password"`
	Email     string         `db:"email" json:"email"`
	Phone     string         `db:"phone" json:"phone"`
	Address   string         `db:"address" json:"address"`
	CreatedAt string         `db:"created_at" json:"created_at"`
	UpdatedAt string         `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullString `db:"deleted_at" json:"deleted_at"`
	IsActive  bool           `db:"is_active" json:"is_active"`
	Role      string         `db:"role" json:"role"`
}
