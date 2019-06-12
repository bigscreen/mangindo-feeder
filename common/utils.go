package common

import (
	"fmt"
	"strings"
)

func GetFormattedChapterNumber(chapter float32) string {
	s := fmt.Sprintf("%.4f", chapter)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}
