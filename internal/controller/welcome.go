package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"net/http"
	"tg_botfather/internal/repository"
	"tg_botfather/internal/services/register"
)

func (h Handler) WelcomeHandler(ctx *gin.Context) {
	update := new(tgbotapi.Update)

	if err := ctx.Bind(update); err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	verified, err := h.Register.CheckUser(ctx, update.Message.From.ID)
	if err != nil {
		if !errors.Is(err, repository.ErrUserNotExists) {
			ctx.Status(http.StatusInternalServerError)
			return
		}
		pleaseVerifyMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Welcome! Please, share your email")
		if _, err := h.Bot.Send(pleaseVerifyMessage); err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}

	if !verified {
		if err := h.Register.ConfirmEmail(ctx, update.Message.Text, update.Message.Chat.ID); err != nil {
			if errors.Is(err, register.ErrBadEmail) {
				badEmailMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Email is wrong, please, repeat")
				if _, err := h.Bot.Send(badEmailMessage); err != nil {
					ctx.Status(http.StatusInternalServerError)
					return
				}
				ctx.Status(http.StatusOK)
				return
			} else {
				ctx.Status(http.StatusInternalServerError)
				return
			}
		}

		emailIsConfirmed := tgbotapi.NewMessage(update.Message.Chat.ID, "Your email is verified")
		if _, err := h.Bot.Send(emailIsConfirmed); err != nil {
			ctx.Status(http.StatusInternalServerError)
			return
		}
	}

	ctx.Status(http.StatusOK)
}
