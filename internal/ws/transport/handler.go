package transport

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/application"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/domain"
)

func (h *Handler) authenticate(r *http.Request) (*auth.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		return h.validateToken(token)
	}

	token := r.Header.Get("Sec-Websocket-Protocol")
	if token != "" {
		return h.validateToken(token)
	}

	return nil, errors.New("missing token")
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	claims, err := h.authenticate(r)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	opts := h.acceptOptions()

	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		slog.Error("ws: failed to accept connection", "error", err)
		return
	}

	client := application.NewClient(claims.Sub)

	defer func() {
		h.manager.UnsubscribeAll(client)
		client.Close()
		_ = conn.CloseNow()
	}()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	writeErr := make(chan error, 1)

	go func() {
		writeErr <- h.writeLoop(ctx, conn, client)
	}()

	if err := h.readLoop(ctx, conn, client); err != nil {
		slog.Debug("ws: read loop ended", "userID", client.UserID, "error", err)
	}

	cancel()
	<-writeErr
}

func (h *Handler) readLoop(ctx context.Context, conn *websocket.Conn, client *application.Client) error {
	for {
		var msg domain.InboundMessage
		if err := wsjson.Read(ctx, conn, &msg); err != nil {
			return err
		}

		if msg.Channel == "" {
			continue
		}

		switch msg.Type {
		case domain.MessageTypeSubscribe:
			h.manager.Subscribe(msg.Channel, client)

		case domain.MessageTypeUnsubscribe:
			h.manager.Unsubscribe(msg.Channel, client)

		case domain.MessageTypeMessage:
			raw, err := json.Marshal(msg.Payload)
			if err != nil {
				slog.Debug("ws: failed to marshal payload", "error", err)
				continue
			}

			var payload any
			if err := json.Unmarshal(raw, &payload); err != nil {
				slog.Debug("ws: failed to unmarshal payload", "error", err)
				continue
			}

			h.manager.Broadcast(msg.Channel, domain.OutboundMessage{
				Channel: msg.Channel,
				From:    client.UserID,
				Payload: payload,
			})

		default:
			slog.Debug("ws: unknown message type", "type", msg.Type)
		}
	}
}

func (h *Handler) writeLoop(ctx context.Context, conn *websocket.Conn, client *application.Client) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg, ok := <-client.Receive():
			if !ok {
				return nil
			}
			if err := wsjson.Write(ctx, conn, msg); err != nil {
				return err
			}
		}
	}
}

func (h *Handler) acceptOptions() *websocket.AcceptOptions {
	for _, o := range h.allowedOrigins {
		if o == "*" {
			return &websocket.AcceptOptions{InsecureSkipVerify: true}
		}
	}

	return &websocket.AcceptOptions{OriginPatterns: h.allowedOrigins}
}
