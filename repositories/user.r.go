package repositories

import (
	"errors"
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
		 RETURNING id, username, email, phone, address, created_at, updated_at, deleted_at, is_active, role`,
		username, email, phone, string(hashedPassword), address,
	).Scan(
		&user.Id,
		&user.Username,
		&user.Email,
		&user.Phone,
		&user.Address,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
		&user.IsActive,
		&user.Role,
	)
	if err != nil {
		log.Printf("DB error (insert): %v", err)
		return errors.New(configs.GetErrString(configs.ErrorCode_FAILED_TO_INSERT_USER))
	}

	return nil
}
