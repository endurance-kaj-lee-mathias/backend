package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type client struct {
	baseURL       string
	realm         string
	adminUser     string
	adminPassword string

	mu          sync.Mutex
	token       string
	tokenExpiry time.Time
}

func NewClient(baseURL string, realm string, adminUser string, adminPassword string) Client {
	return &client{
		baseURL:       baseURL,
		realm:         realm,
		adminUser:     adminUser,
		adminPassword: adminPassword,
	}
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func (c *client) getAdminToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		return c.token, nil
	}

	endpoint := fmt.Sprintf("%s/realms/master/protocol/openid-connect/token", c.baseURL)

	data := url.Values{
		"grant_type": {"password"},
		"client_id":  {"admin-cli"},
		"username":   {c.adminUser},
		"password":   {c.adminPassword},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("keycloak: create token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("keycloak: token request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("keycloak: token request returned %d: %s", resp.StatusCode, body)
	}

	var tok tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tok); err != nil {
		return "", fmt.Errorf("keycloak: decode token: %w", err)
	}

	c.token = tok.AccessToken
	c.tokenExpiry = time.Now().Add(time.Duration(tok.ExpiresIn-30) * time.Second)

	return c.token, nil
}

type keycloakUserPayload struct {
	FirstName  string              `json:"firstName,omitempty"`
	LastName   string              `json:"lastName,omitempty"`
	Email      string              `json:"email,omitempty"`
	Username   string              `json:"username,omitempty"`
	Attributes map[string][]string `json:"attributes,omitempty"`
}

func (c *client) UpdateUser(ctx context.Context, userID string, update UserUpdate) error {
	token, err := c.getAdminToken(ctx)
	if err != nil {
		return err
	}

	payload := keycloakUserPayload{
		FirstName:  update.FirstName,
		LastName:   update.LastName,
		Email:      update.Email,
		Username:   update.Username,
		Attributes: make(map[string][]string),
	}

	if update.PhoneNumber != nil {
		payload.Attributes["phoneNumber"] = []string{*update.PhoneNumber}
	}
	if update.Street != nil {
		payload.Attributes["street"] = []string{*update.Street}
	}
	if update.Locality != nil {
		payload.Attributes["locality"] = []string{*update.Locality}
	}
	if update.Region != nil {
		payload.Attributes["region"] = []string{*update.Region}
	}
	if update.PostalCode != nil {
		payload.Attributes["postal_code"] = []string{*update.PostalCode}
	}
	if update.Country != nil {
		payload.Attributes["country"] = []string{*update.Country}
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("keycloak: marshal payload: %w", err)
	}

	endpoint := fmt.Sprintf("%s/admin/realms/%s/users/%s", c.baseURL, c.realm, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("keycloak: create update request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("keycloak: update request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("keycloak: update user returned %d: %s", resp.StatusCode, respBody)
	}

	return nil
}
