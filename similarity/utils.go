package similarity

import (
	"yasin-wu/utils/common"
)

func GetFinger(fingerArr []string) int64 {
	fingerStr := ""
	for _, f := range fingerArr {
		fingerStr += f
	}
	return common.ConvertString2To10(fingerStr)
}

func GetWords(wordsWeights []WordWeight) []string {
	var words []string
	for _, w := range wordsWeights {
		words = append(words, w.Word)
	}
	return words
}
