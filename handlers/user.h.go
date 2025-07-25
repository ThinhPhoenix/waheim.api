package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"waheim.api/configs"
	"waheim.api/services"
)

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
