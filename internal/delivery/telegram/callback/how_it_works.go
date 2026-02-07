package callback

import (
	"github.com/jus1d/kypidbot/internal/config/messages"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) HowItWorks(c tele.Context) error {
	if err := c.Respond(); err != nil {
		return err
	}
	return c.Send(messages.M.Command.About)
}
