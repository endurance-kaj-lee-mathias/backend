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
	"time"
)

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

func (c *client) getUser(ctx context.Context, token string, userID string) (keycloakUser, error) {
	endpoint := fmt.Sprintf("%s/admin/realms/%s/users/%s", c.baseURL, c.realm, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return keycloakUser{}, fmt.Errorf("keycloak: create get request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return keycloakUser{}, fmt.Errorf("keycloak: get user request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return keycloakUser{}, fmt.Errorf("keycloak: get user returned %d: %s", resp.StatusCode, body)
	}

	var user keycloakUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return keycloakUser{}, fmt.Errorf("keycloak: decode user: %w", err)
	}

	if user.Attributes == nil {
		user.Attributes = make(map[string][]string)
	}

	return user, nil
}

func (c *client) UpdateUser(ctx context.Context, userID string, update UserUpdate) error {
	token, err := c.getAdminToken(ctx)
	if err != nil {
		return err
	}

	existing, err := c.getUser(ctx, token, userID)
	if err != nil {
		return err
	}

	if update.FirstName != "" {
		existing.FirstName = update.FirstName
	}
	if update.LastName != "" {
		existing.LastName = update.LastName
	}
	if update.Email != "" {
		existing.Email = update.Email
	}
	if update.Username != "" {
		existing.Username = update.Username
	}

	if update.PhoneNumber != nil {
		existing.Attributes["phoneNumber"] = []string{*update.PhoneNumber}
	}
	if update.Street != nil {
		existing.Attributes["street"] = []string{*update.Street}
	}
	if update.Locality != nil {
		existing.Attributes["locality"] = []string{*update.Locality}
	}
	if update.Region != nil {
		existing.Attributes["region"] = []string{*update.Region}
	}
	if update.PostalCode != nil {
		existing.Attributes["postal_code"] = []string{*update.PostalCode}
	}
	if update.Country != nil {
		existing.Attributes["country"] = []string{*update.Country}
	}

	body, err := json.Marshal(existing)
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
