//go:build js && wasm

package main

import (
	"syscall/js"
	"unicode"
)

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
	// 期待される戻り値の例：
	//   {"total": 42, "japanese": 30, "noSpace": 38}
	return map[string]int{
		"total":    total,
		"japanese": japanese,
		"noSpace":  noSpace,
	}
}

// jsAnalyzeText is the JS↔Go bridge wrapper registered as a global function.
func jsAnalyzeText(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return js.ValueOf(map[string]interface{}{
			"total":    0,
			"japanese": 0,
			"noSpace":  0,
		})
	}
	text := args[0].String()
	result := analyzeText(text)

	return js.ValueOf(map[string]interface{}{
		"total":    result["total"],
		"japanese": result["japanese"],
		"noSpace":  result["noSpace"],
	})
}

func main() {
	done := make(chan struct{})

	// Register the Go function as a JavaScript global.
	js.Global().Set("analyzeText", js.FuncOf(jsAnalyzeText))

	// Keep the Wasm module alive until the channel receives a value.
	<-done
}

