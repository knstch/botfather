package controller

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"tg_botfather/internal/repository"
	"tg_botfather/internal/services/register"
)

type Handler struct {
	Register register.Register
	Bot      *tgbotapi.BotAPI
}

func NewHandler(db *repository.DbRepository, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{
		Register: *register.NewRegisterService(db),
		Bot:      bot,
	}
}
