package handlers

import (
	"encoding/json"
	"messenger/models"
	request "messenger/requests"
	response "messenger/responses"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

func CreateChatHandler(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var req request.ChatCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request", http.StatusBadRequest)
			return
		}

		title := strings.TrimSpace(req.Title)
		if len(title) < 1 || len(title) > 200 {
			http.Error(w, "title length must be 1..200", http.StatusBadRequest)
			return
		}

		chat := models.Chat{Title: title, CreatedAt: time.Now()}
		if err := db.Create(&chat).Error; err != nil {
			http.Error(w, "db error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chat)
	}
}

func GetChatHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/chats/")
		idStr = strings.Split(idStr, "/")[0]
		chatID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid chat id", http.StatusBadRequest)
			return
		}

		limit := 20
		if l := r.URL.Query().Get("limit"); l != "" {
			if lim, err := strconv.Atoi(l); err == nil && lim > 0 && lim <= 100 {
				limit = lim
			}
		}

		var chat models.Chat
		if err := db.First(&chat, chatID).Error; err != nil {
			http.Error(w, "chat not found", http.StatusNotFound)
			return
		}

		var messages []models.Message
		db.Where("chat_id = ?", chatID).Order("created_at desc").Limit(limit).Find(&messages)
		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}

		resp := response.ChatWithMessages{Chat: chat, Messages: messages}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func DeleteChatHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/chats/")
		chatID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "invalid chat id", http.StatusBadRequest)
			return
		}

		if err := db.Delete(&models.Chat{}, chatID).Error; err != nil {
			http.Error(w, "chat not found", http.StatusNotFound)
			return
		}

		db.Unscoped().Where("chat_id = ?", chatID).Delete(&models.Message{})
		w.WriteHeader(http.StatusNoContent)
	}
}
