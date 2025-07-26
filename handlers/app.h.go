package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"waheim.api/configs"
	"waheim.api/models"
	"waheim.api/services"
)

var appService = services.NewAppService()

func CreateAppHandler(w http.ResponseWriter, r *http.Request) {
	var appReq map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&appReq); err != nil {
		http.Error(w, configs.GetErrString(configs.ErrorCode_INVALID_REQUEST), http.StatusBadRequest)
		return
	}
	// Lấy user_id và role từ context
	userID, _ := r.Context().Value("user_id").(string)
	role, _ := r.Context().Value("role").(string)
	if role != "admin" {
		appReq["publisher_id"] = userID
	}
	// map appReq sang models.App (bạn nên dùng struct, ở đây demo đơn giản)
	// TODO: validate dữ liệu
	// Chuyển map sang struct App (giả sử đã đúng key)
	var app models.App
	if v, ok := appReq["name"].(string); ok {
		app.Name = v
	}
	if v, ok := appReq["description"].(string); ok {
		app.Description = v
	}
	if v, ok := appReq["status"].(string); ok {
		app.Status = v
	}
	if v, ok := appReq["uri"].(string); ok {
		app.Uri = v
	}
	if v, ok := appReq["icon"].(string); ok {
		app.Icon = v
	}
	if v, ok := appReq["publisher_id"].(string); ok {
		app.PublisherId = v
	}
	if v, ok := appReq["category"].(string); ok {
		app.Category = v
	}
	if v, ok := appReq["rating"].(float64); ok {
		app.Rating = v
	}
	if v, ok := appReq["downloads"].(float64); ok {
		app.Downloads = int(v)
	}
	if v, ok := appReq["screenshots"].([]interface{}); ok {
		for _, s := range v {
			if str, ok := s.(string); ok {
				app.ScreenShots = append(app.ScreenShots, str)
			}
		}
	}
	if v, ok := appReq["tags"].([]interface{}); ok {
		for _, s := range v {
			if str, ok := s.(string); ok {
				app.Tags = append(app.Tags, str)
			}
		}
	}
	err := appService.CreateApp(&app)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app)
}

func UpdateAppHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, configs.GetErrString(configs.ErrorCode_MISSING_REQUIRED_FIELDS), http.StatusBadRequest)
		return
	}
	userID, _ := r.Context().Value("user_id").(string)
	role, _ := r.Context().Value("role").(string)
	app, err := appService.GetAppById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if role != "admin" && app.PublisherId != userID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, configs.GetErrString(configs.ErrorCode_INVALID_REQUEST), http.StatusBadRequest)
		return
	}
	err = appService.UpdateApp(id, updates)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteAppHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, configs.GetErrString(configs.ErrorCode_MISSING_REQUIRED_FIELDS), http.StatusBadRequest)
		return
	}
	userID, _ := r.Context().Value("user_id").(string)
	role, _ := r.Context().Value("role").(string)
	app, err := appService.GetAppById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	if role != "admin" && app.PublisherId != userID {
		http.Error(w, "Permission denied", http.StatusForbidden)
		return
	}
	err = appService.DeleteApp(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetAppByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, configs.GetErrString(configs.ErrorCode_MISSING_REQUIRED_FIELDS), http.StatusBadRequest)
		return
	}
	app, err := appService.GetAppById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(app)
}

func GetAllAppsHandler(w http.ResponseWriter, r *http.Request) {
	limit := 10
	offset := 0
	if l := r.URL.Query().Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}
	if o := r.URL.Query().Get("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}
	apps, err := appService.GetAllApps(limit, offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apps)
}
