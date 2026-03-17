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
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/domain"
)

type conversationListerStub struct {
	conversationIDs []uuid.UUID
	err             error
}

func (s conversationListerStub) ListConversationIDs(_ context.Context, _ uuid.UUID) ([]uuid.UUID, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.conversationIDs, nil
}

func TestServeWS_UnauthorizedRequest(t *testing.T) {
	manager := application.NewManager()
	authenticate := func(*http.Request) (*auth.Claims, error) {
		return nil, errors.New("unauthorized")
	}

	handler := NewHandler(manager, conversationListerStub{}, authenticate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	_, _, err := websocket.Dial(context.Background(), wsURL, nil)
	if err == nil {
		t.Fatal("expected unauthorized error, got nil")
	}
}

func TestServeWS_AuthorizedConnection(t *testing.T) {
	manager := application.NewManager()
	userID := uuid.Must(uuid.NewV4())
	authenticate := func(*http.Request) (*auth.Claims, error) {
		return &auth.Claims{Sub: userID.String()}, nil
	}

	handler := NewHandler(manager, conversationListerStub{}, authenticate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func(conn *websocket.Conn) {
		if err := conn.CloseNow(); err != nil {
			t.Logf("failed to close connection: %v", err)
		}
	}(conn)

	if conn == nil {
		t.Fatal("expected connection, got nil")
	}
}

func TestServeWS_AutoSubscribeConversations(t *testing.T) {
	manager := application.NewManager()
	userID := uuid.Must(uuid.NewV4())
	conversationID := uuid.Must(uuid.NewV4())
	authenticate := func(*http.Request) (*auth.Claims, error) {
		return &auth.Claims{Sub: userID.String()}, nil
	}

	lister := conversationListerStub{conversationIDs: []uuid.UUID{conversationID}}
	handler := NewHandler(manager, lister, authenticate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func(conn *websocket.Conn) {
		if err := conn.CloseNow(); err != nil {
			t.Logf("failed to close connection: %v", err)
		}
	}(conn)

	time.Sleep(100 * time.Millisecond)

	channel := "conversation:" + conversationID.String()
	if manager.GetChannelSubscribers(channel) != 1 {
		t.Fatalf("expected 1 subscriber on %s, got %d", channel, manager.GetChannelSubscribers(channel))
	}
}

func TestServeWS_BroadcastToConversationQueue(t *testing.T) {
	manager := application.NewManager()
	userID := uuid.Must(uuid.NewV4())
	conversationID := uuid.Must(uuid.NewV4())
	channel := "conversation:" + conversationID.String()
	authenticate := func(*http.Request) (*auth.Claims, error) {
		return &auth.Claims{Sub: userID.String()}, nil
	}

	lister := conversationListerStub{conversationIDs: []uuid.UUID{conversationID}}
	handler := NewHandler(manager, lister, authenticate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer func(conn *websocket.Conn) {
		if err := conn.CloseNow(); err != nil {
			t.Logf("failed to close connection: %v", err)
		}
	}(conn)

	time.Sleep(100 * time.Millisecond)

	sent := domain.OutboundMessage{
		Channel:   channel,
		SenderID:  userID.String(),
		Content:   "hello",
		CreatedAt: time.Now().UTC(),
	}
	manager.Broadcast(channel, sent)

	var received domain.OutboundMessage
	if err := wsjson.Read(ctx, conn, &received); err != nil {
		t.Fatalf("failed to receive broadcast: %v", err)
	}

	if received.Channel != channel {
		t.Fatalf("expected channel %s, got %s", channel, received.Channel)
	}

	if received.SenderID != userID.String() {
		t.Fatalf("expected sender %s, got %s", userID.String(), received.SenderID)
	}

	if received.Content != "hello" {
		t.Fatalf("expected content hello, got %s", received.Content)
	}
}

func TestServeWS_MultipleClientsSubscribed(t *testing.T) {
	manager := application.NewManager()
	userID := uuid.Must(uuid.NewV4())
	conversationID := uuid.Must(uuid.NewV4())
	channel := "conversation:" + conversationID.String()
	authenticate := func(*http.Request) (*auth.Claims, error) {
		return &auth.Claims{Sub: userID.String()}, nil
	}

	lister := conversationListerStub{conversationIDs: []uuid.UUID{conversationID}}
	handler := NewHandler(manager, lister, authenticate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn1, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect client 1: %v", err)
	}
	defer func(conn *websocket.Conn) {
		if err := conn.CloseNow(); err != nil {
			t.Logf("failed to close client 1: %v", err)
		}
	}(conn1)

	conn2, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect client 2: %v", err)
	}
	defer func(conn *websocket.Conn) {
		if err := conn.CloseNow(); err != nil {
			t.Logf("failed to close client 2: %v", err)
		}
	}(conn2)

	time.Sleep(100 * time.Millisecond)

	if manager.GetChannelSubscribers(channel) != 2 {
		t.Fatalf("expected 2 subscribers, got %d", manager.GetChannelSubscribers(channel))
	}
}

func TestAcceptOptions_WildcardOrigin(t *testing.T) {
	manager := application.NewManager()
	authenticate := func(*http.Request) (*auth.Claims, error) {
		return &auth.Claims{Sub: uuid.Must(uuid.NewV4()).String()}, nil
	}

	handler := NewHandler(manager, conversationListerStub{}, authenticate, []string{"*"})
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
	authenticate := func(*http.Request) (*auth.Claims, error) {
		return &auth.Claims{Sub: uuid.Must(uuid.NewV4()).String()}, nil
	}

	origins := []string{"https://example.com", "https://app.example.com"}
	handler := NewHandler(manager, conversationListerStub{}, authenticate, origins)
	opts := handler.acceptOptions()

	if opts == nil {
		t.Fatal("expected accept options, got nil")
	}
	if len(opts.OriginPatterns) != 2 {
		t.Fatalf("expected 2 origin patterns, got %d", len(opts.OriginPatterns))
	}
	if opts.OriginPatterns[0] != origins[0] {
		t.Fatalf("expected origin %s, got %s", origins[0], opts.OriginPatterns[0])
	}
}

func TestServeWS_ConnectionClosedGracefully(t *testing.T) {
	manager := application.NewManager()
	userID := uuid.Must(uuid.NewV4())
	conversationID := uuid.Must(uuid.NewV4())
	authenticate := func(*http.Request) (*auth.Claims, error) {
		return &auth.Claims{Sub: userID.String()}, nil
	}

	lister := conversationListerStub{conversationIDs: []uuid.UUID{conversationID}}
	handler := NewHandler(manager, lister, authenticate, []string{"*"})
	server := httptest.NewServer(http.HandlerFunc(handler.ServeWS))
	defer server.Close()

	wsURL := "ws" + server.URL[4:]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if err := conn.CloseNow(); err != nil {
		t.Logf("close error (expected): %v", err)
	}

	time.Sleep(100 * time.Millisecond)

	if manager.GetChannels() != 0 {
		t.Fatal("expected all channels to be cleaned up after disconnect")
	}
}
