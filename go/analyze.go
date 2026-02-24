package main

import "unicode"

// analyzeText は指定されたテキストを解析し、文字数の集計結果をマップ形式で返します。
//
// 返り値:
//   - "total":    総文字数（バイト数ではなく、ルーン/文字単位の数）
//   - "japanese": 日本語の文字数（ひらがな ＋ カタカナ ＋ 漢字）
//   - "noSpace":  空白を除いた文字数
func analyzeText(text string) map[string]int {
	total := 0
	japanese := 0
	noSpace := 0

	for _, r := range []rune(text) {
		total++

		if unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r) || unicode.Is(unicode.Han, r) {
			japanese++
		}
		if !unicode.IsSpace(r) {
			noSpace++
		}
	}
	return map[string]int{
		"total":    total,
		"japanese": japanese,
		"noSpace":  noSpace,
	}
}
