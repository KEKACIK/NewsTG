package news

import (
	"newtg/internal/source"
	"time"
)

type NewStatus string

const (
	WaitNewStatus  NewStatus = "wait"
	DoneNewStatus  NewStatus = "done"
	ErrorNewStatus NewStatus = "error"
)

type News struct {
	ID      int           `json:"id"`
	Title   string        `json:"title"`
	Link    string        `json:"link"`
	Content string        `json:"content"`
	Source  source.Source `json:"source"`
	Status  NewStatus     `json:"status"`

	Published time.Time `json:"published"`
	Created   time.Time `json:"created"`
}
