package message

import (
	"fmt"
	"log/slog"

	tele "gopkg.in/telebot.v3"
)

func (h *Handler) Sticker(c tele.Context) error {
	sticker := c.Message().Sticker
	if sticker == nil {
		return nil
	}

	slog.Info("sticker file_id", "file_id", sticker.FileID)
	return c.Send(fmt.Sprintf("```\n%s\n```", sticker.FileID), tele.ModeMarkdown)
}
