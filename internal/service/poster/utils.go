package poster

import (
	"fmt"
	"strings"
)

func (tp *TelegramPoster) ChangeText(text, link string) string {
	if len(strings.Split(text, "")) > tp.maxNewLength {
		text = text[:tp.maxNewLength] + "..."
	}
	text += fmt.Sprintf("\n\n<a href='%s'>Ссылка</a>", link)

	return text
}
