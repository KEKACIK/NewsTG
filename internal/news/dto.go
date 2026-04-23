package news

import "time"

type CreateDTO struct {
	Title     string
	Link      string
	Content   string
	SourceID  int
	Likes     int
	Published time.Time
}

type GetAllDTO struct {
	Status   string
	FromDate time.Time
	ToDate   time.Time
	Limit    int
}
