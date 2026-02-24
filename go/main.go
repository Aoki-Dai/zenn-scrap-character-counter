//go:build js && wasm

package main

import (
	"syscall/js"
)

// jsAnalyzeText は、グローバル関数として登録された JS ↔ Go 間のブリッジ用ラッパー
func jsAnalyzeText(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(map[string]any{
			"total":    0,
			"japanese": 0,
			"noSpace":  0,
		})
	}
	text := args[0].String()
	result := analyzeText(text)

	return js.ValueOf(map[string]any{
		"total":    result["total"],
		"japanese": result["japanese"],
		"noSpace":  result["noSpace"],
	})
}

func main() {
	done := make(chan struct{})

	// Goの関数をJavaScriptのグローバル関数として登録
	js.Global().Set("analyzeText", js.FuncOf(jsAnalyzeText))

	// チャネルが値を受信するまで、Wasmモジュールを維持
	<-done
}
