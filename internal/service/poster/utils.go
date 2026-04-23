package poster

import (
	"fmt"
	"strings"
)

func (tp *TelegramPoster) ChangeText(text, link string) string {
	if len(strings.Split(text, "")) > tp.maxNewLength {
		text = string([]rune(text)[:tp.maxNewLength]) + "..."
	}

	text = strings.Join(
		[]string{
			text,
			"",
			fmt.Sprintf("<a href='%s'>Источник</a>", link),
		},
		"\n",
	)

	return text
}
