package main

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDecodeTelegramResponseUnwrapsSuccessfulResult(t *testing.T) {
	var user TelegramUser

	err := decodeTelegramResponse([]byte(`{"ok":true,"result":{"id":123,"is_bot":true,"first_name":"Mistle","username":"mistle_bot"}}`), &user)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != 123 || !user.IsBot || user.FirstName != "Mistle" || user.Username != "mistle_bot" {
		t.Fatalf("unexpected user: %#v", user)
	}
}

func TestDecodeTelegramResponseReturnsProviderError(t *testing.T) {
	var out json.RawMessage

	err := decodeTelegramResponse([]byte(`{"ok":false,"error_code":400,"description":"Bad Request: chat not found"}`), &out)
	if err == nil {
		t.Fatal("expected provider error")
	}

	if !strings.Contains(err.Error(), "telegram api error 400: Bad Request: chat not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}
