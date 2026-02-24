package main

import "testing"

func TestAnalyzeText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		total    int
		japanese int
		noSpace  int
	}{
		{
			name:     "空文字列",
			input:    "",
			total:    0,
			japanese: 0,
			noSpace:  0,
		},
		{
			name:     "英数字のみ",
			input:    "Hello, World!",
			total:    13,
			japanese: 0,
			noSpace:  12,
		},
		{
			name:     "日本語のみ",
			input:    "こんにちは、世界！",
			total:    9,
			japanese: 7,
			noSpace:  9,
		},
		{
			name:     "英数字と日本語",
			input:    "Hello, 世界!",
			total:    10,
			japanese: 2,
			noSpace:  9,
		},
		{
			name:     "スペースのみ",
			input:    "   ",
			total:    3,
			japanese: 0,
			noSpace:  0,
		},
		{
			name:     "スペースと日本語",
			input:    "  こんにちは、世界！",
			total:    11,
			japanese: 7,
			noSpace:  9,
		},
		{
			name:     "スペースと英数字",
			input:    "  Hello, World!",
			total:    15,
			japanese: 0,
			noSpace:  12,
		},
		{
			name:     "スペースと英数字と日本語",
			input:    "  Hello, 世界!",
			total:    12,
			japanese: 2,
			noSpace:  9,
		},
		{
			name:     "スペースと英数字と日本語とスペース",
			input:    "  Hello, 世界!  ",
			total:    14,
			japanese: 2,
			noSpace:  9,
		},
		{
			name:  "カタカナ",
			input: "ハローワールド",
			total: 7,
			// ー（U+30FC）は unicode.Katakana の範囲（U+30A1–U+30FA, U+30FD–U+30FF）に含まれないため japanese=5
			japanese: 5,
			noSpace:  7,
		},
		// 全角スペース（U+3000）は unicode.IsSpace で space 扱いされる
		{
			name:     "全角スペースのみ",
			input:    "　",
			total:    1,
			japanese: 0,
			noSpace:  0,
		},
		{
			name:     "全角スペースと日本語",
			input:    "　こんにちは",
			total:    6,
			japanese: 5,
			noSpace:  5,
		},
		// 改行は unicode.IsSpace で space 扱いされる（実際のスクラップ本文を想定）
		{
			name:     "改行を含む日本語",
			input:    "こんにちは\n世界",
			total:    8,
			japanese: 7,
			noSpace:  7,
		},
		// 半角カタカナ（U+FF66–FF9D）は unicode.Katakana の範囲に含まれるため japanese にカウントされる
		{
			name:     "半角カタカナ",
			input:    "ｱｲｳ",
			total:    3,
			japanese: 3,
			noSpace:  3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := analyzeText(tt.input)
			if got["total"] != tt.total {
				t.Errorf("total: got %d, want %d", got["total"], tt.total)
			}
			if got["japanese"] != tt.japanese {
				t.Errorf("japanese: got %d, want %d", got["japanese"], tt.japanese)
			}
			if got["noSpace"] != tt.noSpace {
				t.Errorf("noSpace: got %d, want %d", got["noSpace"], tt.noSpace)
			}
		})
	}
}
