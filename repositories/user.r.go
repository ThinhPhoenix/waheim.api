package repositories

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
	"waheim.api/configs"
	"waheim.api/models"
)

func SignIn(request map[string]string) (string, error) {
	db := configs.DB
	var user models.User

	waheimId, hasId := request["waheim_id"]
	password, hasPass := request["password"]
	if !hasId || !hasPass || waheimId == "" || password == "" {
		return "", errors.New(configs.GetErrString(configs.ErrorCode_SIGN_IN_MISSING_FIELDS))
	}

	query := `SELECT * FROM users WHERE username ILIKE $1 OR email ILIKE $1 OR phone ILIKE $1 LIMIT 1`
	err := db.Get(&user, query, waheimId)
	if err != nil {
		log.Printf("SignIn DB error: %v, waheim_id: %s", err, waheimId)
		return "", errors.New(configs.GetErrString(configs.ErrorCode_AUTH_FAILED))
	}

	if !user.IsActive {
		return "", errors.New(configs.GetErrString(configs.ErrorCode_USER_NOT_ACTIVE))
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New(configs.GetErrString(configs.ErrorCode_AUTH_FAILED))
	}

	tokenString, err := configs.GenerateJwt(user.Id, user.Role)
	if err != nil {
		return "", errors.New(configs.GetErrString(configs.ErrorCode_FAILED_TO_GENERATE_TOKEN))
	}

	return tokenString, nil
}

func AuthMe(tokenString string) (models.User, error) {
	claims, err := configs.ValidateJwt(tokenString)
	if err != nil {
		return models.User{}, errors.New(configs.GetErrString(configs.ErrorCode_INVALID_TOKEN))
	}
	userId, ok := claims["user_id"].(string)
	if !ok {
		return models.User{}, errors.New(configs.GetErrString(configs.ErrorCode_INVALID_USER_ID_IN_TOKEN))
	}
	db := configs.DB
	var user models.User
	err = db.Get(&user, "SELECT * FROM users WHERE id = $1 LIMIT 1", userId)
	if err != nil {
		return user, errors.New(configs.GetErrString(configs.ErrorCode_USER_NOT_FOUND))
	}
	return user, nil
}

func SignUp(request map[string]string) error {
	db := configs.DB

	username := request["username"]
	email := request["email"]
	phone := request["phone"]
	password := request["password"]
	address := request["address"]

	if username == "" || email == "" || phone == "" || password == "" {
		return errors.New(configs.GetErrString(configs.ErrorCode_SIGN_UP_MISSING_FIELDS))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Password hash error: %v", err)
		return errors.New(configs.GetErrString(configs.ErrorCode_FAILED_TO_HASH_PASSWORD))
	}

	var exists int
	err = db.Get(&exists, "SELECT COUNT(*) FROM users WHERE username=$1 OR email=$2 OR phone=$3", username, email, phone)
	if err != nil {
		log.Printf("DB error (check exists): %v", err)
		return errors.New(configs.GetErrString(configs.ErrorCode_DATABASE_ERROR))
	}
	if exists > 0 {
		return errors.New(configs.GetErrString(configs.ErrorCode_USER_ALREADY_EXISTS))
	}

	var user models.User
	err = db.QueryRowx(
		`INSERT INTO users (username, email, phone, password, address, is_active, role, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, true, 'user', NOW(), NOW())
		 RETURNING id, username, password, email, phone, address, created_at, updated_at, deleted_at, is_active, role, avatar, first_name, last_name, date_of_birth, gender, status`,
		username, email, phone, string(hashedPassword), address,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.IsActive,
		&user.Role,
		&user.Avatar,
		&user.FirstName,
		&user.LastName,
		&user.DateOfBirth,
		&user.Gender,
		&user.Status,
	)
	if err != nil {
		log.Printf("DB error (insert): %v", err)
		return errors.New(configs.GetErrString(configs.ErrorCode_FAILED_TO_INSERT_USER))
	}

	return nil
}

func GetAllUsers(filters map[string]string, limit, offset int) ([]models.User, error) {
	db := configs.DB
	var users []models.User
	query := "SELECT * FROM users WHERE deleted_at IS NULL"
	args := []interface{}{}
	idx := 1

	for key, raw := range filters {
		if raw == "" {
			continue
		}
		op := "="
		val := raw
		if len(raw) > 2 && (raw[:2] == ">=" || raw[:2] == "<=") {
			op = raw[:2]
			val = raw[2:]
		} else if len(raw) > 1 && (raw[:1] == ">" || raw[:1] == "<") {
			op = raw[:1]
			val = raw[1:]
		} else if len(raw) > 1 && raw[:1] == "=" {
			op = "="
			val = raw[1:]
		} else if len(raw) > 1 && raw[:1] == "%" {
			op = "LIKE"
		}
		if op == "LIKE" {
			query += fmt.Sprintf(" AND %s ILIKE $%d", key, idx)
			args = append(args, val)
		} else {
			query += fmt.Sprintf(" AND %s %s $%d", key, op, idx)
			args = append(args, val)
		}
		idx++
	}

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", idx)
		args = append(args, limit)
		idx++
	}
	if offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", idx)
		args = append(args, offset)
		idx++
	}

	err := db.Select(&users, query, args...)
	if err != nil {
		log.Printf("DB error (get all users): %v", err)
		return nil, errors.New(configs.GetErrString(configs.ErrorCode_DATABASE_ERROR))
	}
	return users, nil
}

// Lấy user theo id
func GetUserById(id string) (models.User, error) {
	db := configs.DB
	var user models.User
	err := db.Get(&user, "SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL", id)
	if err != nil {
		log.Printf("DB error (get user by id): %v", err)
		return user, errors.New(configs.GetErrString(configs.ErrorCode_USER_NOT_FOUND))
	}
	return user, nil
}

func UpdateUser(id string, updates map[string]interface{}) error {
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
	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d AND deleted_at IS NULL", setClause, idx)
	res, err := db.Exec(query, args...)
	if err != nil {
		log.Printf("DB error (update user): %v", err)
		return errors.New(configs.GetErrString(configs.ErrorCode_DATABASE_ERROR))
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(configs.GetErrString(configs.ErrorCode_USER_NOT_FOUND))
	}
	return nil
}

func DeleteUser(id string) error {
	db := configs.DB
	query := "UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL"
	res, err := db.Exec(query, id)
	if err != nil {
		log.Printf("DB error (delete user): %v", err)
		return errors.New(configs.GetErrString(configs.ErrorCode_DATABASE_ERROR))
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return errors.New(configs.GetErrString(configs.ErrorCode_USER_NOT_FOUND))
	}
	return nil
}
