package configs

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"os"
	"io"
)

func SendToCloud(filePath string) (*http.Response, error) {
	apiUrl := os.Getenv("TELEREALM_URI")
	botToken := os.Getenv("TELEREALM_BOT_TOKEN")
	chatID := os.Getenv("TELEREALM_CHAT_ID")

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	_ = writer.WriteField("chat_id", chatID)

	part, err := writer.CreateFormFile("document", file.Name())
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", apiUrl, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+botToken)

	client := &http.Client{}
	return client.Do(req)
}