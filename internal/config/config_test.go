package config

import (
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	cfg := Config{
		Model:         "translategemma",
		OpenAIBaseURL: "http://localhost:11434/v1",
		OpenAIKey:     "sk-test",
		MaxConcurrent: 10,
	}

	if cfg.Model != "translategemma" {
		t.Errorf("Model = %q, want %q", cfg.Model, "translategemma")
	}
	if cfg.OpenAIBaseURL != "http://localhost:11434/v1" {
		t.Errorf("OpenAIBaseURL = %q, want %q", cfg.OpenAIBaseURL, "http://localhost:11434/v1")
	}
	if cfg.OpenAIKey != "sk-test" {
		t.Errorf("OpenAIKey = %q, want %q", cfg.OpenAIKey, "sk-test")
	}
	if cfg.MaxConcurrent != 10 {
		t.Errorf("MaxConcurrent = %d, want %d", cfg.MaxConcurrent, 10)
	}
}
