package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"waheim.api/configs"
	"waheim.api/services"
)

// Adapter cho Gin -> http.Handler, truyền param vào context
type paramsKeyType struct{}

var paramsKey = paramsKeyType{}

func GinToHTTPHandler(h http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := map[string]string{}
		for _, p := range c.Params {
			params[p.Key] = p.Value
		}
		ctx := context.WithValue(c.Request.Context(), paramsKey, params)
		h(c.Writer, c.Request.WithContext(ctx))
	}
}

var userService = services.NewUserService()

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, configs.GetErrString(configs.ErrorCode_INVALID_REQUEST), http.StatusBadRequest)
		return
	}
	err := userService.SignUp(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	var req map[string]string
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, configs.GetErrString(configs.ErrorCode_INVALID_REQUEST), http.StatusBadRequest)
		return
	}
	token, err := userService.SignIn(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	secure := r.URL.Scheme == "https" || r.TLS != nil
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   r.URL.Scheme == "https" || r.TLS != nil,
		SameSite: http.SameSiteLaxMode,
	})
	w.WriteHeader(http.StatusOK)
}

func AuthMeHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, configs.GetErrString(configs.ErrorCode_MISSING_AUTH_HEADER), http.StatusUnauthorized)
		return
	}
	var token string
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		token = authHeader[7:]
	} else {
		http.Error(w, configs.GetErrString(configs.ErrorCode_INVALID_AUTH_HEADER_FORMAT), http.StatusUnauthorized)
		return
	}
	resp, err := userService.AuthMe(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	filters := make(map[string]string)
	limit := 10
	offset := 0

	if username := r.URL.Query().Get("username"); username != "" {
		filters["username"] = username
	}
	if email := r.URL.Query().Get("email"); email != "" {
		filters["email"] = email
	}
	if phone := r.URL.Query().Get("phone"); phone != "" {
		filters["phone"] = phone
	}
	if isActive := r.URL.Query().Get("is_active"); isActive != "" {
		filters["is_active"] = isActive
	}
	if createdAt := r.URL.Query().Get("created_at"); createdAt != "" {
		filters["created_at"] = createdAt
	}
	if role := r.URL.Query().Get("role"); role != "" {
		filters["role"] = role
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil {
			limit = parsedLimit
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil {
			offset = parsedOffset
		}
	}

	users, err := userService.GetAllUsers(filters, limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := ""
	if r.Context().Value(paramsKey) != nil {
		if params, ok := r.Context().Value(paramsKey).(map[string]string); ok {
			id = params["id"]
		}
	}
	if id == "" {
		id = r.URL.Query().Get("id")
	}
	if id == "" {
		http.Error(w, configs.GetErrString(configs.ErrorCode_MISSING_REQUIRED_FIELDS), http.StatusBadRequest)
		return
	}
	// Lấy user_id và role từ context
	userID, _ := r.Context().Value("user_id").(string)
	role, _ := r.Context().Value("role").(string)
	if role != "admin" && userID != id {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	err := userService.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := ""
	if r.Context().Value(paramsKey) != nil {
		if params, ok := r.Context().Value(paramsKey).(map[string]string); ok {
			id = params["id"]
		}
	}
	if id == "" {
		id = r.URL.Query().Get("id")
	}
	if id == "" {
		http.Error(w, configs.GetErrString(configs.ErrorCode_MISSING_REQUIRED_FIELDS), http.StatusBadRequest)
		return
	}
	user, err := userService.GetUserById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := ""
	if r.Context().Value(paramsKey) != nil {
		if params, ok := r.Context().Value(paramsKey).(map[string]string); ok {
			id = params["id"]
		}
	}
	if id == "" {
		id = r.URL.Query().Get("id")
	}
	if id == "" {
		http.Error(w, configs.GetErrString(configs.ErrorCode_MISSING_REQUIRED_FIELDS), http.StatusBadRequest)
		return
	}
	// Lấy user_id và role từ context
	userID, _ := r.Context().Value("user_id").(string)
	role, _ := r.Context().Value("role").(string)
	if role != "admin" && userID != id {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, configs.GetErrString(configs.ErrorCode_INVALID_REQUEST), http.StatusBadRequest)
		return
	}
	err := userService.UpdateUser(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
