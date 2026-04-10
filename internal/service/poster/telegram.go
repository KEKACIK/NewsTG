package poster

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramPoster struct {
	bot    *tgbotapi.BotAPI
	chatID int64
}

func NewTelegramPoster(token string, chatID int64) (*TelegramPoster, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to init bot: %w", err)
	}

	return &TelegramPoster{
		bot:    bot,
		chatID: chatID,
	}, nil
}

func (tp *TelegramPoster) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(tp.chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := tp.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
