package parser

import (
	"context"
	"fmt"
	"net/http"
	"newtg/internal/news"
	"newtg/internal/source"
	"newtg/pkg/logging"
	"newtg/pkg/postgresql"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type RiaClient struct {
	client postgresql.Client
	logger *logging.Logger

	sourceName string
	tags       []string
	url        string
	dynamicUrl string

	dataFormat     string
	timeFormat     string
	maxNewsPerHour int
}

type RiaPost struct {
	Title       string
	Link        string
	Content     string
	LikeSum     int
	PublishedAt time.Time
}

func NewRiaClient(
	client postgresql.Client,
	logger *logging.Logger,

	sourceName string,
	maxNewsPerHour int,
) *RiaClient {
	return &RiaClient{
		client: client,
		logger: logger,

		sourceName: sourceName,
		tags: []string{
			"world", "incidents", "economy", "society",
			"politics", "culture", "tourism", "science",
			"defense_safety", "religion", "sport",
		},
		url:            "https://ria.ru/services/tagsearch/?date_start=%s&date=%s&tags[]=%s",
		dynamicUrl:     "https://ria.ru/services/dynamics/%s/%s.html",
		dataFormat:     "20060102",
		timeFormat:     "15:04 02.01.2006",
		maxNewsPerHour: maxNewsPerHour,
	}
}

func (rc *RiaClient) PoolNews(ctx context.Context) {
	dataFrom := time.Now().AddDate(0, 0, -1)
	dataTo := time.Now()

	allPosts := make([]*RiaPost, 0)
	linksCheck := map[string]struct{}{}
	for _, tag := range rc.tags {
		posts, err := rc.GetAllPosts(linksCheck, dataFrom, dataTo, tag)
		if err != nil {
			rc.logger.Warn(fmt.Sprintf("Ошибка чтения HTML: %s", err))
			continue
		}
		allPosts = append(allPosts, posts...)
	}

	sourceRepo := source.NewRepository(rc.client, rc.logger)
	rSource, err := sourceRepo.GetByName(context.Background(), rc.sourceName)
	if err != nil {
		rc.logger.Fatal(err.Error())
	}

	newsRepo := news.NewRepository(rc.client, rc.logger)

	for _, post := range allPosts {
		err = newsRepo.Create(context.Background(), &news.CreateDTO{
			Title:     post.Title,
			Link:      post.Link,
			Content:   post.Content,
			SourceID:  rSource.ID,
			Likes:     post.LikeSum,
			Published: post.PublishedAt,
		})
		if err != nil {
			rc.logger.Error(err.Error())
		}
	}
}

func (rc *RiaClient) GetAllPosts(linkCheck map[string]struct{}, dataFrom, dataTo time.Time, tag string) ([]*RiaPost, error) {
	hClient := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf(rc.url, dataFrom.Format(rc.dataFormat), dataTo.Format(rc.dataFormat), tag)
	res, err := hClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Ошибка: статус код %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	posts := make([]*RiaPost, 0)

	doc.Find(".list-item").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Find("a.list-item__title.color-font-hover-only").Attr("href")
		if !exists {
			rc.logger.Warn(err.Error())
			return
		}

		if _, exists := linkCheck[link]; exists {
			return
		}

		post, err := rc.GetPost(hClient, link)
		if err != nil {
			rc.logger.Warn(err.Error())
			return
		}

		linkCheck[link] = struct{}{}
		posts = append(posts, post)
	})

	return posts, nil
}

func (rc *RiaClient) GetPost(hClient *http.Client, postUrl string) (*RiaPost, error) {
	res, err := hClient.Get(postUrl)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Ошибка: статус код %d", res.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	post := &RiaPost{}

	post.Title = doc.Find(".article__title").First().Text()

	post.Link = postUrl
	{
		contentList := make([]string, 0)
		doc.Find(".article__summary-list li").Each(func(i int, s *goquery.Selection) {
			contentList = append(contentList, s.Text())
		})
		if len(contentList) < 1 {
			contentList = append(contentList, doc.Find(".article__text").First().Text())
		}
		post.Content = strings.Join(contentList, "\n\n")
	}

	{
		dynamicUrl, err := rc.getDynamicLink(postUrl)
		if err != nil {
			return nil, err
		}

		dynamicRes, err := hClient.Get(dynamicUrl)
		if err != nil {
			return nil, err
		}
		defer dynamicRes.Body.Close()

		if dynamicRes.StatusCode != 200 {
			return nil, fmt.Errorf("Ошибка: статус код %d", dynamicRes.StatusCode)
		}

		dynamicDoc, err := goquery.NewDocumentFromReader(dynamicRes.Body)
		if err != nil {
			return nil, err
		}

		for i := 1; i <= 6; i++ {
			countStr := dynamicDoc.Find(fmt.Sprintf(".emoji-item.m-type-s%d .m-value", i)).First().Text()
			count, err := strconv.Atoi(countStr)
			if err != nil {
				return nil, err
			}

			post.LikeSum += count
		}
	}

	post.PublishedAt, err = time.Parse(rc.timeFormat, doc.Find("div.article__info-date a").First().Text())
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (rc *RiaClient) getDynamicLink(link string) (string, error) {
	re := regexp.MustCompile(`ria\.ru/(\d{8})/.+-(\d+)\.html`)
	matches := re.FindStringSubmatch(link)

	if len(matches) < 3 {
		fmt.Println("Не удалось распарсить ссылку")
		return "", fmt.Errorf("Dynamic link error")
	}

	return fmt.Sprintf(rc.dynamicUrl, matches[1], matches[2]), nil
}
