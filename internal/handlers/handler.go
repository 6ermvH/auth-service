package handlers

import (
	"encoding/json"
	"net"
	"net/http"
	"log"

	"example.com/auth_service/repository"
	"example.com/auth_service/mail"
	"example.com/auth_service/token"
)

type Handler struct {
	Repo repository.IRefreshTokenRepository
	TokenM token.ITokenManager
}

type IHandler interface {
	HandleGenerateTokens(http.ResponseWriter, *http.Request)
	HandleUpdateTokens(http.ResponseWriter, *http.Request)
}

func NewHandler (repo repository.IRefreshTokenRepository, secretKey string) *Handler {
	if len(secretKey) == 0 {
		log.Fatalf("Can`t find JWT_SECRET_KEY")
		return nil
	}
	return &Handler{Repo: repo, TokenM: token.NewTokenManager([]byte(secretKey))}
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

	accessToken, refreshToken, err := h.TokenM.GenerateNew(userID, clientIP)
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

	parsedUserID, err := h.TokenM.ParseAccess(req.AccessToken, "user_id")
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusBadRequest)
		return
	}

	parsedClientIP, err := h.TokenM.ParseAccess(req.AccessToken, "ip")
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusBadRequest)
		return
	}

	if clientIP != parsedClientIP {
		err = mail.SendToIP(clientIP, "Warning. Someone has access to your account.")
		if err != nil {
			log.Panicf("Failed to send email to %s", clientIP)
		}
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

	accessToken, refreshToken, err := h.TokenM.GenerateNew(parsedUserID, clientIP)
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
