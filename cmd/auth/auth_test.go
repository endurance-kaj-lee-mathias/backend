package auth

import (
	"net/http/httptest"
	"testing"
)

func TestExtractToken_FromAuthorizationHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer abc.def.ghi")

	token, err := extractToken(req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if token != "abc.def.ghi" {
		t.Fatalf("expected token abc.def.ghi, got %s", token)
	}
}

func TestExtractToken_FromWebsocketProtocolHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Sec-Websocket-Protocol", "abc.def.ghi")

	token, err := extractToken(req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if token != "abc.def.ghi" {
		t.Fatalf("expected token abc.def.ghi, got %s", token)
	}
}

func TestExtractToken_AuthorizationTakesPrecedence(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer fromAuthorization")
	req.Header.Set("Sec-Websocket-Protocol", "fromProtocol")

	token, err := extractToken(req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if token != "fromAuthorization" {
		t.Fatalf("expected token fromAuthorization, got %s", token)
	}
}

func TestExtractToken_MissingHeaders(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)

	_, err := extractToken(req)
	if err != MissingHeader {
		t.Fatalf("expected MissingHeader, got %v", err)
	}
}
