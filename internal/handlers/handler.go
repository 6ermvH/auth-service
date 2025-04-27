package handlers

import (
	"encoding/json"
	"net"
	"net/http"

	"example.com/auth_service/token"
	"example.com/auth_service/repository"
)

type Handler struct {
	Repo repository.IRefreshTokenRepository
}

type IHandler interface {
	HandleGenerateTokens(http.ResponseWriter, *http.Request)
	HandleUpdateTokens(http.ResponseWriter, *http.Request)
}

func (h *Handler) HandleGenerateTokens(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Invalid client IP address", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := token.GenerateNew(userID, clientIP)
	if err != nil {
		http.Error(w, "Failed to generate tokens", http.StatusInternalServerError)
		return
	}

	err = h.Repo.InsertRefreshToken(userID, clientIP, refreshToken)
	if err != nil {
		http.Error(w, "Failed to insert token to DataBase", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleUpdateTokens(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	accessToken, refreshToken, err := token.Update(req.AccessToken, req.RefreshToken)
	if err != nil {
		http.Error(w, "Failed to update tokens", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
