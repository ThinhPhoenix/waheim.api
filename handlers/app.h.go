package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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

	// Gọi GitHub API để build APK
	githubToken := os.Getenv("GITHUB_TOKEN")
	repoOwner := os.Getenv("GITHUB_REPO_OWNER")
	repoName := os.Getenv("GITHUB_REPO_NAME")
	apiUrl := fmt.Sprintf("https://api.github.com/repos/%s/%s/dispatches", repoOwner, repoName)
	apkUrl := ""
	// Dùng app.Uri làm url web build APK
	payload := fmt.Sprintf(`{"event_type":"build-apk","client_payload":{"url":"%s"}}`, app.Uri)
	req, _ := http.NewRequest("POST", apiUrl, bytes.NewBuffer([]byte(payload)))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+githubToken)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode >= 300 {
		// Không build được thì vẫn trả app, chỉ log lỗi
		if err == nil {
			body, _ := ioutil.ReadAll(resp.Body)
			err = fmt.Errorf("github dispatch failed: %s", string(body))
		}
		fmt.Println("Error triggering GitHub build:", err)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(app)
		return
	}
	// Đợi workflow build xong và lấy artifact (polling đơn giản, thực tế nên dùng async)
	runId := ""
	for i := 0; i < 20; i++ { // poll tối đa 20 lần, mỗi lần 15s
		time.Sleep(15 * time.Second)
		// Lấy workflow runs
		runsReq, _ := http.NewRequest("GET",
			fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs?event=repository_dispatch", repoOwner, repoName), nil)
		runsReq.Header.Set("Authorization", "Bearer "+githubToken)
		runsResp, err := client.Do(runsReq)
		if err != nil {
			continue
		}
		runsBody, _ := ioutil.ReadAll(runsResp.Body)
		runsStr := string(runsBody)
		// Tìm runId mới nhất
		idx := strings.Index(runsStr, "\"id\":")
		if idx > 0 {
			runsStr = runsStr[idx+6:]
			end := strings.IndexAny(runsStr, ",}")
			if end > 0 {
				runId = runsStr[:end]
			}
		}
		if runId != "" {
			// Kiểm tra run đã success chưa
			checkReq, _ := http.NewRequest("GET",
				fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%s", repoOwner, repoName, runId), nil)
			checkReq.Header.Set("Authorization", "Bearer "+githubToken)
			checkResp, err := client.Do(checkReq)
			if err != nil {
				continue
			}
			checkBody, _ := ioutil.ReadAll(checkResp.Body)
			if strings.Contains(string(checkBody), "\"conclusion\":\"success\"") {
				// Lấy artifact
				artReq, _ := http.NewRequest("GET",
					fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%s/artifacts", repoOwner, repoName, runId), nil)
				artReq.Header.Set("Authorization", "Bearer "+githubToken)
				artResp, err := client.Do(artReq)
				if err != nil {
					break
				}
				artBody, _ := ioutil.ReadAll(artResp.Body)
				// Tìm url download
				idx := strings.Index(string(artBody), "archive_download_url")
				if idx > 0 {
					artBodyStr := string(artBody)[idx+23:]
					end := strings.IndexAny(artBodyStr, "\",}")
					if end > 0 {
						apkUrl = artBodyStr[:end]
					}
				}
				break
			}
		}
	}
	if apkUrl != "" {
		// Update AndroidInstallUri cho app
		updates := map[string]interface{}{"android_install_uri": apkUrl}
		_ = appService.UpdateApp(app.Id, updates)
		app.AndroidInstallUri = apkUrl
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
