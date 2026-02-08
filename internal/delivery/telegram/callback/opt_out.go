package callback

import (
	"context"
	"log/slog"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/delivery/telegram/view"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) OptOut(c tele.Context) error {
	sender := c.Sender()

	user, err := h.Users.GetUser(context.Background(), sender.ID)
	if err != nil || user == nil {
		slog.Error("get user for opt out", sl.Err(err))
		return c.Respond()
	}

	optedout := !user.OptedOut
	if err := h.Users.SetOptedOut(context.Background(), sender.ID, optedout); err != nil {
		slog.Error("set opted out", sl.Err(err))
		return c.Respond()
	}

	_ = c.Respond()
	return c.Edit(messages.M.Registration.Completed, view.RegistrationCompletedKeyboard(optedout))
}
