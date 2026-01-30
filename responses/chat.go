package response

import "messenger/models"

type ChatWithMessages struct {
	Chat     models.Chat      `json:"chat"`
	Messages []models.Message `json:"messages"`
}
