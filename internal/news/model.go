package news

import (
	"newtg/internal/source"
	"time"
)

type News struct {
	ID        int           `json:"id"`
	Title     string        `json:"title"`
	Link      string        `json:"link"`
	Content   string        `json:"content"`
	Source    source.Source `json:"source"`
	Posted    bool          `json:"posted"`
	Published time.Time     `json:"published"`
	Created   time.Time     `json:"created"`
}
