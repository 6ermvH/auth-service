package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"log"

	"example.com/auth_service/repository"
	"example.com/auth_service/token"
)

type Handler struct {
	Repo repository.IRefreshTokenRepository
}

type IHandler interface {
	HandleGenerateTokens(http.ResponseWriter, *http.Request)
	HandleUpdateTokens(http.ResponseWriter, *http.Request)
}

func (h *Handler) HandleGenerateTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid http method", http.StatusBadRequest)
		return
	}

	log.Printf("Handle generate is active. CLIENT: %s\nWITH DATA: %s", r.RemoteAddr, r.Body)

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
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid http method", http.StatusBadRequest)
		return
	}

	log.Printf("Handle generate is active. CLIENT: %s\nWITH DATA: %s", r.RemoteAddr, r.Body)
	var req struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)

	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	parsedUserID, err := token.ParseAccess(req.AccessToken, "user_id")
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusBadRequest)
		return
	}

	parsedClientIP, err := token.ParseAccess(req.AccessToken, "ip")
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusBadRequest)
		return
	}

	if clientIP != parsedClientIP {
		http.Error(w, "Other ip", http.StatusBadRequest)
		return
	}

	ok, err := h.Repo.CheckRefreshToken(parsedUserID, req.RefreshToken)
	if err != nil {
		http.Error(w, "Failed to check refresh token", http.StatusInternalServerError)
		return
	} else if !ok {
		http.Error(w, "Invalid refresh token", http.StatusBadRequest)
		return
	}

	err = h.Repo.MarkRefreshTokenUsed(parsedUserID, req.RefreshToken)
	if err != nil {
		http.Error(w, "Failed to update refresh token", http.StatusInternalServerError)
		return
	}

	accessToken, refreshToken, err := token.GenerateNew(parsedUserID, clientIP)
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
