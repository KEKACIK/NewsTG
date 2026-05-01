package poster

import (
	"context"
	"errors"
	"fmt"
	"newtg/internal/news"
	"newtg/internal/source"
	"newtg/pkg/logging"
	"newtg/pkg/postgresql"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/robfig/cron/v3"
	"gopkg.in/telebot.v4"
)

type TelegramPoster struct {
	client postgresql.Client
	logger *logging.Logger
	// Telegram
	bot    *telebot.Bot
	chatID int64
	// Ria news
	riaSourceName        string
	riaLimitPerHour      int
	telegramMaxMsgLength int
}

func NewTelegramPoster(
	client postgresql.Client,
	logger *logging.Logger,
	// Telegram
	token string,
	chatID int64,
	telegramMaxMsgLength int,
	// Ria news
	riaSourceName string,
	riaLimitPerHour int,
) (*TelegramPoster, error) {
	bot_settings := telebot.Settings{Token: token}
	bot, err := telebot.NewBot(bot_settings)
	if err != nil {
		return nil, err
	}

	return &TelegramPoster{
		client:               client,
		logger:               logger,
		bot:                  bot,
		chatID:               chatID,
		riaSourceName:        riaSourceName,
		riaLimitPerHour:      riaLimitPerHour,
		telegramMaxMsgLength: telegramMaxMsgLength,
	}, nil
}

func (tp *TelegramPoster) StartPool(ctx context.Context) error {
	sourceRepo := source.NewRepository(tp.client, tp.logger)

	riaSource, err := sourceRepo.GetByName(ctx, tp.riaSourceName)
	if err != nil {
		return err
	}

	c := cron.New()
	c.AddFunc("00 * * * *", func() {
		tp.logger.Info("Start posting")
		go tp.CheckRiaNews(context.Background(), riaSource.ID, tp.riaLimitPerHour)
	})
	c.Start()

	tp.logger.Info("Poster start pooling..")
	select {}
}

func (tp *TelegramPoster) CheckRiaNews(ctx context.Context, source_id int, limit int) {
	newsRepo := news.NewRepository(tp.client, tp.logger)

	waitNews, err := newsRepo.GetAll(ctx, &news.GetAllDTO{
		SourceID: source_id,
		Status:   string(news.WaitStatus),
		FromDate: time.Now().Add(-(2 * time.Hour)),
		Limit:    limit,
	})
	if err != nil {
		tp.logger.Error(err.Error())
		return
	}
	if len(waitNews) == 0 {
		tp.logger.Info("No news to send")
		return
	}

	NewsToSend := make([]news.News, 0)
	TextToSend := make([]string, 0)
	for _, newItem := range waitNews {
		tp.logger.Debug(fmt.Sprintf("Загружаю news.%d...", newItem.ID))

		title := newItem.Title
		runesTitle := []rune(newItem.Title)
		if len(runesTitle) > 100 {
			title = string(runesTitle[:100])
		}

		content := newItem.Content
		runesContent := []rune(content)
		if len(runesContent) > 1100 {
			truncated := string(runesContent[:1100])
			lastDot := strings.LastIndex(truncated, ".")
			if lastDot > 0 {
				content = truncated[:lastDot+1]
			}
		}

		TextToSend = append(TextToSend, strings.Join([]string{
			fmt.Sprintf("<b>%s</b>", title),
			"",
			fmt.Sprintf("<blockquote>%s</blockquote>", content),
			"",
			fmt.Sprintf("<a href='%s'>%s</a>", newItem.Link, tp.riaSourceName),
		}, "\n"))
		NewsToSend = append(NewsToSend, newItem)
	}

	text := strings.Join(TextToSend, "\n━━━━━━━━━━━\n")

	err = tp.BotSendPost(ctx, text)
	if err != nil {
		tp.logger.Warn(err.Error())
		return
	}

	for _, newItem := range NewsToSend {
		newItem.Status = news.DoneStatus
		err = newsRepo.Update(ctx, &newItem)
		if err != nil {
			tp.logger.Warn(err.Error())
		}
	}
}

func (tp *TelegramPoster) BotSendPost(ctx context.Context, text string) error {
	text = "🤖 <b>Новости часа</b>\n\n" + text
	if utf8.RuneCountInString(text) > tp.telegramMaxMsgLength {
		return errors.New("Post too long, skipping")
	}

	_, err := tp.bot.Send(
		&telebot.Chat{ID: tp.chatID},
		text,
		&telebot.SendOptions{
			ParseMode:             telebot.ModeHTML,
			DisableWebPagePreview: true,
		},
	)

	return err
}
