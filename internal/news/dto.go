package news

type CreateNewsDTO struct {
	Title     string `json:"title"`
	Link      string `json:"link"`
	Content   string `json:"content"`
	SourceID  int    `json:"source_id"`
	Published string `json:"published"`
	Posted    string `json:"posted"`
}
