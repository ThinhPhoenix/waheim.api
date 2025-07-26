package models

type Rating struct {
	Id        string `db:"id" json:"id"`
	UserId    string `db:"user_id" json:"user_id"`
	AppId     string `db:"app_id" json:"app_id"`
	Comment   string `db:"comment" json:"comment"`
	Stars     int    `db:"stars" json:"stars"`
	CreatedAt string `db:"created_at" json:"created_at"`
	UpdatedAt string `db:"updated_at" json:"updated_at"`
	DeletedAt string `db:"deleted_at" json:"deleted_at"`
	Status    string `db:"status" json:"status"`
}
