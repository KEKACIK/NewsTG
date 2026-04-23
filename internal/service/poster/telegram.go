package poster

import (
	"context"
	"fmt"
	"newtg/internal/news"
	"newtg/pkg/logging"
	"newtg/pkg/postgresql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramPoster struct {
	client postgresql.Client
	logger *logging.Logger

	bot    *tgbotapi.BotAPI
	chatID int64

	maxNewLength int
}

func (tp *TelegramPoster) CheckNews(ctx context.Context) {
	newsRepo := news.NewRepository(tp.client, tp.logger)
	waitNews, err := newsRepo.GetAllByStatus(context.TODO(), news.WaitNewStatus)
	if err != nil {
		tp.logger.Error(err.Error())
		return
	}

	for _, wn := range waitNews {
		tp.logger.Debug(fmt.Sprintf("Обработка news.%d начата...", wn.ID))

		text := tp.ChangeText(wn.Content, wn.Link)

		err = tp.SendMessage(text)
		wn.Status = news.DoneNewStatus
		if err != nil {
			wn.Status = news.ErrorNewStatus
			tp.logger.Warn(fmt.Sprintf("news.%d статус '%s'. Error: %s", wn.ID, wn.Status, err.Error()))
		}

		err = newsRepo.Update(ctx, &wn)
		if err != nil {
			tp.logger.Error(err.Error())
		}
	}
}

func (tp *TelegramPoster) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(tp.chatID, text)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = true

	_, err := tp.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func NewTelegramPoster(
	client postgresql.Client,
	logger *logging.Logger,

	token string,
	chatID int64,

	maxNewLength int,
) (*TelegramPoster, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to init bot: %w", err)
	}

	return &TelegramPoster{
		client: client,
		logger: logger,

		bot:    bot,
		chatID: chatID,

		maxNewLength: maxNewLength,
	}, nil
}
