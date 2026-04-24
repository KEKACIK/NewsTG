package poster

import (
	"context"
	"fmt"
	"newtg/internal/news"
	"newtg/internal/source"
	"newtg/pkg/logging"
	"newtg/pkg/postgresql"
	"strings"
	"time"

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
	riaSourceName   string
	riaLimitPerHour int
}

func (tp *TelegramPoster) StartPool(ctx context.Context) error {
	sourceRepo := source.NewRepository(tp.client, tp.logger)

	riaSource, err := sourceRepo.GetByName(ctx, tp.riaSourceName)
	if err != nil {
		return err
	}

	c := cron.New()
	c.AddFunc("00 * * * *", func() {
		go tp.CheckRiaNews(context.Background(), riaSource.ID, tp.riaLimitPerHour)
	})
	c.Start()

	tp.logger.Info("Poster start pooling..")
	select {}
}

func (tp *TelegramPoster) CheckRiaNews(ctx context.Context, source_id int, limit int) {
	newsRepo := news.NewRepository(tp.client, tp.logger)

	waitNews, err := newsRepo.GetAll(context.Background(), &news.GetAllDTO{
		Status:   string(news.WaitStatus),
		FromDate: time.Now().Add(-time.Hour),
		Limit:    limit,
	})
	if err != nil {
		tp.logger.Error(err.Error())
		return
	}

	for _, wn := range waitNews {
		tp.logger.Debug(fmt.Sprintf("Отправка news.%d...", wn.ID))
		wn.Status = news.DoneStatus

		err = tp.BotSendPost(ctx, tp.riaSourceName, &wn)
		if err != nil {
			wn.Status = news.ErrorStatus
			tp.logger.Warn(fmt.Sprintf("news.%d статус '%s'. Error: %s", wn.ID, wn.Status, err.Error()))
		}

		err = newsRepo.Update(ctx, &wn)
		if err != nil {
			tp.logger.Warn(err.Error())
		}
		time.Sleep(5 * time.Second)
	}
}

func (tp *TelegramPoster) BotSendPost(ctx context.Context, linkName string, new *news.News) error {
	recipient := telebot.Chat{ID: tp.chatID}
	text := strings.Join([]string{
		fmt.Sprintf("<b>%s</b>", new.Title),
		"",
		fmt.Sprintf("<blockquote>%s</blockquote>", new.Content),
		"",
		fmt.Sprintf("<a href='%s'>%s</a>", new.Link, linkName),
	}, "\n")
	options := telebot.SendOptions{
		ParseMode:             telebot.ModeHTML,
		DisableWebPagePreview: true,
	}

	_, err := tp.bot.Send(
		&recipient,
		text,
		&options,
	)

	return err
}

func NewTelegramPoster(
	client postgresql.Client,
	logger *logging.Logger,
	// Telegram
	token string,
	chatID int64,
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
		client:          client,
		logger:          logger,
		bot:             bot,
		chatID:          chatID,
		riaSourceName:   riaSourceName,
		riaLimitPerHour: riaLimitPerHour,
	}, nil
}
