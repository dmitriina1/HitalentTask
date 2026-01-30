package tests

import (
	"bytes"
	"encoding/json"
	"messenger/handlers"
	"messenger/models"
	request "messenger/requests"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Chat{}, &models.Message{})
	return db
}

func TestCreateChatHandler(t *testing.T) {
	db := setupTestDB()
	reqBody := request.ChatCreateRequest{Title: "Test Chat"}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/chats/", bytes.NewReader(b))
	w := httptest.NewRecorder()
	handlers.CreateChatHandler(db)(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var chat models.Chat
	json.NewDecoder(w.Body).Decode(&chat)
	assert.Equal(t, "Test Chat", chat.Title)
}

func TestGetChatHandler(t *testing.T) {
	db := setupTestDB()
	chat := models.Chat{Title: "Chat"}
	db.Create(&chat)
	for i := 0; i < 3; i++ {
		db.Create(&models.Message{ChatID: chat.ID, Text: "Msg"})
	}
	req := httptest.NewRequest(http.MethodGet, "/chats/1?limit=2", nil)
	w := httptest.NewRecorder()
	handlers.GetChatHandler(db)(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteChatHandler(t *testing.T) {
	db := setupTestDB()
	chat := models.Chat{Title: "Chat"}
	db.Create(&chat)
	db.Create(&models.Message{ChatID: chat.ID, Text: "Msg"})
	req := httptest.NewRequest(http.MethodDelete, "/chats/1", nil)
	w := httptest.NewRecorder()
	handlers.DeleteChatHandler(db)(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
	var count int64
	db.Model(&models.Message{}).Count(&count)
	assert.Equal(t, int64(0), count)
}
