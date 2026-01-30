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
)

func TestSendMessageHandler(t *testing.T) {
	db := setupTestDB()
	chat := models.Chat{Title: "Chat"}
	db.Create(&chat)
	reqBody := request.MessageCreateRequest{Text: "Hello, world!"}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/chats/1/messages/", bytes.NewReader(b))
	w := httptest.NewRecorder()
	handlers.SendMessageHandler(db)(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var msg models.Message
	json.NewDecoder(w.Body).Decode(&msg)
	assert.Equal(t, "Hello, world!", msg.Text)
}
