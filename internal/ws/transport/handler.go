package transport

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"github.com/go-chi/chi/v5"
	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/cmd/auth"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/response"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/application"
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

	userID, err := uuid.FromString(claims.Sub)
	if err != nil {
		response.WriteError(w, http.StatusUnauthorized, Unauthorized)
		return
	}

	conversationID, err := uuid.FromString(chi.URLParam(r, "conversationId"))
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, InvalidConversationID)
		return
	}

	conversationIDs, err := h.conversations.ListConversationIDs(r.Context(), userID)
	if err != nil {
		slog.Error("ws: failed to list conversation subscriptions", "userID", claims.Sub, "error", err)
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	if !containsConversation(conversationIDs, conversationID) {
		response.WriteError(w, http.StatusForbidden, Forbidden)
		return
	}

	opts := h.acceptOptions()

	conn, err := websocket.Accept(w, r, opts)
	if err != nil {
		slog.Error("ws: failed to accept connection", "error", err)
		return
	}

	client := application.NewClient(claims.Sub)
	h.manager.Subscribe(conversationChannel(conversationID), client)

	defer func() {
		h.manager.UnsubscribeAll(client)
		client.Close()
		_ = conn.CloseNow()
	}()

	ctx := conn.CloseRead(r.Context())

	if err := h.writeLoop(ctx, conn, client); err != nil && err != context.Canceled {
		slog.Debug("ws: write loop ended", "userID", client.UserID, "error", err)
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

func conversationChannel(conversationID uuid.UUID) string {
	return "conversation:" + conversationID.String()
}

func containsConversation(conversationIDs []uuid.UUID, conversationID uuid.UUID) bool {
	for _, id := range conversationIDs {
		if id == conversationID {
			return true
		}
	}

	return false
}