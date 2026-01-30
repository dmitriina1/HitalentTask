package handlers

import (
	"encoding/json"
	"messenger/models"
	request "messenger/requests"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func SendMessageHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/chats/")
		idStr = strings.Split(idStr, "/")[0]
		chatID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid chat id", http.StatusBadRequest)
			return
		}

		var req request.MessageCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		text := strings.TrimSpace(req.Text)
		if len(text) < 1 || len(text) > 5000 {
			http.Error(w, "text length must be 1..5000", http.StatusBadRequest)
			return
		}

		var chat models.Chat
		if err := db.First(&chat, chatID).Error; err != nil {
			http.Error(w, "chat not found", http.StatusNotFound)
			return
		}

		msg := models.Message{ChatID: chatID, Text: text, CreatedAt: time.Now()}
		if err := db.Create(&msg).Error; err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)
	}
}
