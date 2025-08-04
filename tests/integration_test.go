package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

const (
	baseURL = "http://localhost:8080"
	userID  = "11111111-1111-1111-1111-111111111111"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func TestGenerateTokens(t *testing.T) {
	url := fmt.Sprintf("%s/token?user_id=%s", baseURL, userID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	var tokenResponse TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	if tokenResponse.AccessToken == "" {
		t.Error("access token is empty")
	}
	if tokenResponse.RefreshToken == "" {
		t.Error("refresh token is empty")
	}
}

func TestRefreshTokens(t *testing.T) {
	// 1. Generate initial tokens
	generateURL := fmt.Sprintf("%s/token?user_id=%s", baseURL, userID)
	req, err := http.NewRequest("POST", generateURL, nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("could not send request: %v", err)
	}
	defer resp.Body.Close()

	var initialTokens TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&initialTokens); err != nil {
		t.Fatalf("could not decode initial response: %v", err)
	}

	// 2. Use the tokens to refresh
	refreshURL := fmt.Sprintf("%s/refresh", baseURL)
	refreshBody, _ := json.Marshal(map[string]string{
		"access_token":  initialTokens.AccessToken,
		"refresh_token": initialTokens.RefreshToken,
	})

	req, err = http.NewRequest("POST", refreshURL, bytes.NewBuffer(refreshBody))
	if err != nil {
		t.Fatalf("could not create refresh request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("could not send refresh request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK for refresh; got %v", resp.Status)
	}

	var refreshedTokens TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&refreshedTokens); err != nil {
		t.Fatalf("could not decode refreshed response: %v", err)
	}

	if refreshedTokens.AccessToken == "" {
		t.Error("refreshed access token is empty")
	}
	if refreshedTokens.RefreshToken == "" {
		t.Error("refreshed refresh token is empty")
	}
}
