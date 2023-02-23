package similarity

import (
	"github.com/yasin-wu/utils/strings"
)

func GetFinger(fingerArr []string) int64 {
	fingerStr := ""
	for _, f := range fingerArr {
		fingerStr += f
	}
	return strings.Base10(fingerStr)
}

func GetWords(wordsWeights []WordWeight) []string {
	var words []string //nolint:prealloc
	for _, w := range wordsWeights {
		words = append(words, w.Word)
	}
	return words
}
