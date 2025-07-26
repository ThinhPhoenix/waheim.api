package models

import (
	"database/sql"

	"github.com/lib/pq"
)

type App struct {
	Id                string         `db:"id" json:"id"`
	Name              string         `db:"name" json:"name"`
	Description       string         `db:"description" json:"description"`
	CreatedAt         string         `db:"created_at" json:"created_at"`
	UpdatedAt         string         `db:"updated_at" json:"updated_at"`
	DeletedAt         sql.NullString `db:"deleted_at" json:"deleted_at"`
	Status            string         `db:"status" json:"status"`
	Uri               string         `db:"uri" json:"uri"`
	Icon              string         `db:"icon" json:"icon"`
	PublisherId       string         `db:"publisher_id" json:"publisher_id"`
	ScreenShots       pq.StringArray `db:"screenshots" json:"screenshots"`
	Category          string         `db:"category" json:"category"`
	Tags              pq.StringArray `db:"tags" json:"tags"`
	Rating            float64        `db:"rating" json:"rating"`
	Downloads         int            `db:"downloads" json:"downloads"`
	AndroidInstallUri string         `db:"android_install_uri" json:"android_install_uri"`
	IOSInstallUri     string         `db:"ios_install_uri" json:"ios_install_uri"`
}
