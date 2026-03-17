package transport

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/domain"
)

func TestServeWS_UnauthorizedRequest(t *testing.T) {
	manager := application.NewManager()
	validate := func(string) (*auth.Claims, error) {
		return nil, errors.New("unauthorized")
	}

	handler := NewHandler(manager, validate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	_, _, err := websocket.Dial(context.Background(), wsURL, dialOptionsAuth("Bearer token"))
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}
}

func TestServeWS_AuthorizedConnection(t *testing.T) {
	manager := application.NewManager()
	validate := func(token string) (*auth.Claims, error) {
		if token == "token123" {
			return &auth.Claims{Sub: "user123"}, nil
		}

		return nil, errors.New("unauthorized")
	}

	handler := NewHandler(manager, validate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, dialOptionsAuth("Bearer token123"))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.CloseNow()
		if err != nil {
			derr := "failed to close connection: %v"
			t.Logf(derr, err)
		}
	}(conn)

	if conn == nil {
		t.Fatal("expected connection, got nil")
	}
}

func TestServeWS_AuthorizedWithSecWebsocketProtocol(t *testing.T) {
	manager := application.NewManager()
	validate := func(token string) (*auth.Claims, error) {
		if token == "token456" {
			return &auth.Claims{Sub: "user456"}, nil
		}

		return nil, errors.New("unauthorized")
	}

	handler := NewHandler(manager, validate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, dialOptionsProtocol("token456"))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.CloseNow()
		if err != nil {
			derr := "failed to close connection: %v"
			t.Logf(derr, err)
		}
	}(conn)

	if conn == nil {
		t.Fatal("expected connection, got nil")
	}
}

func TestServeWS_Subscribe(t *testing.T) {
	manager := application.NewManager()
	validate := func(string) (*auth.Claims, error) {
		return &auth.Claims{Sub: "user123"}, nil
	}

	handler := NewHandler(manager, validate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, dialOptionsAuth("Bearer token"))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.CloseNow()
		if err != nil {
			derr := "failed to close connection: %v"
			t.Logf(derr, err)
		}
	}(conn)

	subscribeMsg := domain.InboundMessage{
		Type:    domain.MessageTypeSubscribe,
		Channel: "general",
	}

	if err := wsjson.Write(ctx, conn, subscribeMsg); err != nil {
		t.Fatalf("failed to send subscribe message: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if manager.GetChannelSubscribers("general") == 0 {
		t.Fatal("expected client to be subscribed to channel")
	}
}

func TestServeWS_Unsubscribe(t *testing.T) {
	manager := application.NewManager()
	validate := func(string) (*auth.Claims, error) {
		return &auth.Claims{Sub: "user123"}, nil
	}

	handler := NewHandler(manager, validate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, dialOptionsAuth("Bearer token"))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.CloseNow()
		if err != nil {
			derr := "failed to close connection: %v"
			t.Logf(derr, err)
		}
	}(conn)

	subscribeMsg := domain.InboundMessage{
		Type:    domain.MessageTypeSubscribe,
		Channel: "general",
	}
	if err := wsjson.Write(ctx, conn, subscribeMsg); err != nil {
		t.Fatalf("failed to send subscribe message: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if manager.GetChannelSubscribers("general") != 1 {
		t.Fatal("expected client to be subscribed")
	}

	unsubscribeMsg := domain.InboundMessage{
		Type:    domain.MessageTypeUnsubscribe,
		Channel: "general",
	}
	if err := wsjson.Write(ctx, conn, unsubscribeMsg); err != nil {
		t.Fatalf("failed to send unsubscribe message: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if manager.GetChannelSubscribers("general") != 0 {
		t.Fatal("expected client to be unsubscribed")
	}
}

func TestServeWS_Broadcast(t *testing.T) {
	manager := application.NewManager()
	validate := func(string) (*auth.Claims, error) {
		return &auth.Claims{Sub: "user123"}, nil
	}

	handler := NewHandler(manager, validate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, dialOptionsAuth("Bearer token"))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.CloseNow()
		if err != nil {
			derr := "failed to close connection: %v"
			t.Logf(derr, err)
		}
	}(conn)

	subscribeMsg := domain.InboundMessage{
		Type:    domain.MessageTypeSubscribe,
		Channel: "general",
	}
	if err := wsjson.Write(ctx, conn, subscribeMsg); err != nil {
		t.Fatalf("failed to send subscribe message: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	broadcastMsg := domain.InboundMessage{
		Type:    domain.MessageTypeMessage,
		Channel: "general",
		Payload: map[string]string{"text": "hello"},
	}
	if err := wsjson.Write(ctx, conn, broadcastMsg); err != nil {
		t.Fatalf("failed to send message: %v", err)
	}

	var received domain.OutboundMessage
	if err := wsjson.Read(ctx, conn, &received); err != nil {
		t.Fatalf("failed to receive message: %v", err)
	}

	if received.Channel != "general" {
		t.Fatalf("expected channel 'general', got '%s'", received.Channel)
	}
	if received.From != "user123" {
		t.Fatalf("expected from 'user123', got '%s'", received.From)
	}
}

func TestServeWS_MultipleClients(t *testing.T) {
	manager := application.NewManager()
	validate := func(string) (*auth.Claims, error) {
		return &auth.Claims{Sub: "user123"}, nil
	}

	handler := NewHandler(manager, validate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn1, _, err := websocket.Dial(ctx, wsURL, dialOptionsAuth("Bearer token"))
	if err != nil {
		t.Fatalf("failed to connect client 1: %v", err)
	}
	defer func(conn1 *websocket.Conn) {
		err := conn1.CloseNow()
		if err != nil {
			derr := "failed to close connection: %v"
			t.Logf(derr, err)
		}
	}(conn1)

	conn2, _, err := websocket.Dial(ctx, wsURL, dialOptionsAuth("Bearer token"))
	if err != nil {
		t.Fatalf("failed to connect client 2: %v", err)
	}
	defer func(conn2 *websocket.Conn) {
		err := conn2.CloseNow()
		if err != nil {
			derr := "failed to close connection: %v"
			t.Logf(derr, err)
		}
	}(conn2)

	subscribeMsg := domain.InboundMessage{
		Type:    domain.MessageTypeSubscribe,
		Channel: "general",
	}

	if err := wsjson.Write(ctx, conn1, subscribeMsg); err != nil {
		t.Fatalf("failed to send subscribe message from client 1: %v", err)
	}

	if err := wsjson.Write(ctx, conn2, subscribeMsg); err != nil {
		t.Fatalf("failed to send subscribe message from client 2: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if manager.GetChannelSubscribers("general") != 2 {
		t.Fatalf("expected 2 clients subscribed, got %d", manager.GetChannelSubscribers("general"))
	}
}

func TestServeWS_EmptyChannelIgnored(t *testing.T) {
	manager := application.NewManager()
	validate := func(string) (*auth.Claims, error) {
		return &auth.Claims{Sub: "user123"}, nil
	}

	handler := NewHandler(manager, validate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, dialOptionsAuth("Bearer token"))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func(conn *websocket.Conn) {
		err := conn.CloseNow()
		if err != nil {
			derr := "failed to close connection: %v"
			t.Logf(derr, err)
		}
	}(conn)

	emptyChannelMsg := domain.InboundMessage{
		Type:    domain.MessageTypeSubscribe,
		Channel: "",
	}

	if err := wsjson.Write(ctx, conn, emptyChannelMsg); err != nil {
		t.Fatalf("failed to send message with empty channel: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if manager.GetChannels() != 0 {
		t.Fatal("expected no channels to be created for empty channel message")
	}
}

func TestAcceptOptions_WildcardOrigin(t *testing.T) {
	manager := application.NewManager()
	validate := func(string) (*auth.Claims, error) {
		return &auth.Claims{Sub: "user123"}, nil
	}

	handler := NewHandler(manager, validate, []string{"*"})
	opts := handler.acceptOptions()

	if opts == nil {
		t.Fatal("expected accept options, got nil")
	}
	if !opts.InsecureSkipVerify {
		t.Fatal("expected InsecureSkipVerify to be true for wildcard origin")
	}
}

func TestAcceptOptions_SpecificOrigins(t *testing.T) {
	manager := application.NewManager()
	validate := func(string) (*auth.Claims, error) {
		return &auth.Claims{Sub: "user123"}, nil
	}

	origins := []string{"https://example.com", "https://app.example.com"}
	handler := NewHandler(manager, validate, origins)
	opts := handler.acceptOptions()

	if opts == nil {
		t.Fatal("expected accept options, got nil")
	}
	if len(opts.OriginPatterns) != 2 {
		t.Fatalf("expected 2 origin patterns, got %d", len(opts.OriginPatterns))
	}
	if opts.OriginPatterns[0] != origins[0] {
		t.Fatalf("expected origin '%s', got '%s'", origins[0], opts.OriginPatterns[0])
	}
}

func TestServeWS_ConnectionClosedGracefully(t *testing.T) {
	manager := application.NewManager()
	validate := func(string) (*auth.Claims, error) {
		return &auth.Claims{Sub: "user123"}, nil
	}

	handler := NewHandler(manager, validate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, dialOptionsAuth("Bearer token"))
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	err = conn.CloseNow()
	if err != nil {
		t.Logf("close error (expected): %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if manager.GetChannels() != 0 {
		t.Fatal("expected all channels to be cleaned up after disconnect")
	}
}

func dialOptionsAuth(value string) *websocket.DialOptions {
	header := http.Header{}
	header.Set("Authorization", value)

	return &websocket.DialOptions{HTTPHeader: header}
}

func dialOptionsProtocol(value string) *websocket.DialOptions {
	header := http.Header{}
	header.Set("Sec-Websocket-Protocol", value)

	return &websocket.DialOptions{HTTPHeader: header}
}
