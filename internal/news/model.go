package news

import "newtg/internal/source"

type News struct {
	ID        int           `json:"id"`
	Title     string        `json:"title"`
	Link      string        `json:"link"`
	Content   string        `json:"content"`
	Source    source.Source `json:"source"`
	Published string        `json:"published"`
	Created   string        `json:"created"`
	Posted    string        `json:"posted"`
}
