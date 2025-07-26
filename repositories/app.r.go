package repositories

import (
	"errors"
	"fmt"
	"log"

	"github.com/lib/pq"
	"waheim.api/configs"
	"waheim.api/models"
)

func CreateApp(app *models.App) error {
	db := configs.DB
	query := `INSERT INTO apps (name, description, status, uri, icon, publisher_id, screenshots, category, tags, rating, downloads, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,NOW(),NOW())
		RETURNING id, created_at, updated_at, deleted_at`
	return db.QueryRowx(query,
		app.Name,
		app.Description,
		app.Status,
		app.Uri,
		app.Icon,
		app.PublisherId,
		pq.StringArray(app.ScreenShots),
		app.Category,
		pq.StringArray(app.Tags),
		app.Rating,
		app.Downloads,
	).Scan(&app.Id, &app.CreatedAt, &app.UpdatedAt, &app.DeletedAt)
}

func GetAppById(id string) (models.App, error) {
	db := configs.DB
	var app models.App
	err := db.Get(&app, "SELECT * FROM apps WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		log.Printf("DB error (get app by id): %v", err)
		return app, errors.New(configs.GetErrString(configs.ErrorCode_APP_NOT_FOUND))
	}
	return app, nil
}

func GetAllApps(limit, offset int) ([]models.App, error) {
	db := configs.DB
	var apps []models.App
	query := "SELECT * FROM apps WHERE deleted_at IS NULL"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", offset)
	}
	err := db.Select(&apps, query)
	if err != nil {
		log.Printf("DB error (get all apps): %v", err)
		return nil, errors.New(configs.GetErrString(configs.ErrorCode_DATABASE_ERROR))
	}
	return apps, nil
}

func UpdateApp(id string, updates map[string]interface{}) error {
	db := configs.DB
	if len(updates) == 0 {
		return nil
	}
	setClause := ""
	args := []interface{}{}
	idx := 1
	for k, v := range updates {
		if setClause != "" {
			setClause += ", "
		}
		setClause += fmt.Sprintf("%s = $%d", k, idx)
		args = append(args, v)
		idx++
	}
	setClause += ", updated_at = NOW()"
	args = append(args, id)
	query := fmt.Sprintf("UPDATE apps SET %s WHERE id = $%d AND deleted_at IS NULL", setClause, idx)
	res, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("DB error (update app): %v", err)
		return errors.New(configs.GetErrString(configs.ErrorCode_DATABASE_ERROR))
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(configs.GetErrString(configs.ErrorCode_APP_NOT_FOUND))
	}
	return nil
}

func DeleteApp(id string) error {
	db := configs.DB
	query := "UPDATE apps SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL"
	res, err := db.Exec(query, id)
	if err != nil {
		log.Printf("DB error (delete app): %v", err)
		return errors.New(configs.GetErrString(configs.ErrorCode_DATABASE_ERROR))
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(configs.GetErrString(configs.ErrorCode_APP_NOT_FOUND))
	}
	return nil
}
