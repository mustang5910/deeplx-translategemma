package logic

import (
	"testing"

	"github.com/mustang5910/deeplx-translategemma/internal/types"
)

func TestResolveLanguageCode(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedCode string
		expectedLang string
	}{
		{
			name:         "known EN code",
			input:        "EN",
			expectedCode: "en",
			expectedLang: "English",
		},
		{
			name:         "known ZH code",
			input:        "ZH",
			expectedCode: "zh",
			expectedLang: "Chinese",
		},
		{
			name:         "empty input",
			input:        "",
			expectedCode: "",
			expectedLang: "",
		},
		{
			name:         "unknown code falls back to input",
			input:        "XX",
			expectedCode: "XX",
			expectedLang: "XX",
		},
		{
			name:         "region fallback EN-US via base EN",
			input:        "EN-US",
			expectedCode: "en",
			expectedLang: "English",
		},
		{
			name:         "ZH-HANS direct map exists",
			input:        "ZH-HANS",
			expectedCode: "zh-Hans",
			expectedLang: "Chinese",
		},
		{
			name:         "lowercase input uppercased before lookup",
			input:        "ja",
			expectedCode: "ja",
			expectedLang: "Japanese",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, lang := resolveLanguageCode(tt.input)
			if code != tt.expectedCode {
				t.Errorf("resolveLanguageCode(%q) code = %q, want %q", tt.input, code, tt.expectedCode)
			}
			if lang != tt.expectedLang {
				t.Errorf("resolveLanguageCode(%q) lang = %q, want %q", tt.input, lang, tt.expectedLang)
			}
		})
	}
}

func TestBuildTranslateParams(t *testing.T) {
	req := &types.Request{
		Text:       "Hello",
		SourceLang: "EN",
		TargetLang: "ZH",
	}

	params := buildTranslateParams(req)

	if params.Text != "Hello" {
		t.Errorf("Text = %q, want %q", params.Text, "Hello")
	}
	if params.SourceCode != "en" {
		t.Errorf("SourceCode = %q, want %q", params.SourceCode, "en")
	}
	if params.SourceLang != "English" {
		t.Errorf("SourceLang = %q, want %q", params.SourceLang, "English")
	}
	if params.TargetCode != "zh" {
		t.Errorf("TargetCode = %q, want %q", params.TargetCode, "zh")
	}
	if params.TargetLang != "Chinese" {
		t.Errorf("TargetLang = %q, want %q", params.TargetLang, "Chinese")
	}
}
