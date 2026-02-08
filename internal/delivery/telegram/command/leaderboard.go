package command

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jus1d/kypidbot/internal/config/messages"
	"github.com/jus1d/kypidbot/internal/lib/logger/sl"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) Leaderboard(c tele.Context) error {
    sender := c.Sender()
    
    slog.Info("processing leaderboard command", 
        "user_id", sender.ID, 
        "username", sender.Username)
    
    leaderboard, err := h.Registration.GetReferralLeaderboard(context.Background(), 10)
    if err != nil {
        slog.Error("failed to get referral leaderboard", 
            sl.Err(err), 
            "user_id", sender.ID)
        
        return c.Send("Произошла ошибка при получении лидерборда 😔", tele.ModeMarkdown)
    }
    
    if len(leaderboard) == 0 {
        return c.Send(messages.M.Command.Leaderboard.Empty, tele.ModeMarkdown)
    }
    
    var messageBuilder strings.Builder
    
    // Добавляем заголовок
    messageBuilder.WriteString(messages.M.Command.Leaderboard.Title) // ← ИЗМЕНИЛИ
    messageBuilder.WriteString("\n\n")
    
    // Добавляем места с эмодзи
    for i, entry := range leaderboard {
        // Определяем эмодзи для места
        var emoji string
        switch i {
        case 0:
            emoji = "🥇"
        case 1:
            emoji = "🥇"
        case 2:
            emoji = "🥇"
        default:
            emoji = fmt.Sprintf("%d.", i+1)
        }
        
        // Определяем имя пользователя для отображения
        displayName := entry.FirstName
        if entry.Username != "" {
            displayName = "@" + entry.Username
        } else if entry.FirstName == "" {
            displayName = fmt.Sprintf("ID: %d", entry.ReferrerID)
        }
        
        // Форматируем строку
        messageBuilder.WriteString(fmt.Sprintf(
            "%s %s — *%d* %s\n",
            emoji,
            displayName,
            entry.ReferralCount,
            h.pluralizeReferrals(entry.ReferralCount),
        ))
    }
    
    // Добавляем разделитель и информацию о текущем пользователе
    messageBuilder.WriteString("\n" + messages.M.Command.Leaderboard.Footer) // ← ИЗМЕНИЛИ
    
    // Отправляем сообщение
    return c.Send(messageBuilder.String(), tele.ModeMarkdown)
}

// Вспомогательная функция для правильного склонения слова "реферал"
func (h *Handler) pluralizeReferrals(count int) string {
    lastDigit := count % 10
    lastTwoDigits := count % 100
    
    if lastTwoDigits >= 11 && lastTwoDigits <= 14 {
        return "рефералов"
    }
    
    switch lastDigit {
    case 1:
        return "реферал"
    case 2, 3, 4:
        return "реферала"
    default:
        return "рефералов"
    }
}