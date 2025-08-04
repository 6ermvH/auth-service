package handlers

import (
	"encoding/json"
	"log"
	"net"
	"net/http"

	"example.com/auth_service/mail"
	"example.com/auth_service/repository"
	"example.com/auth_service/token"
)

type Handler struct {
	Repo   repository.IRefreshTokenRepository
	TokenM token.ITokenManager
}

type IHandler interface {
	HandleGenerateTokens(http.ResponseWriter, *http.Request)
	HandleUpdateTokens(http.ResponseWriter, *http.Request)
}

func NewHandler(repo repository.IRefreshTokenRepository, secretKey string) *Handler {
	if len(secretKey) == 0 {
		log.Fatalf("Can't find JWT_SECRET_KEY")
		return nil
	}
	return &Handler{Repo: repo, TokenM: token.NewTokenManager([]byte(secretKey))}
}

func (h *Handler) HandleGenerateTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusBadRequest, "Invalid http method")
		return
	}

	log.Printf("Handle generate is active. CLIENT: %s", r.RemoteAddr)

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		respondWithError(w, http.StatusBadRequest, "Missing user_id")
		return
	}

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid client IP address")
		return
	}

	accessToken, refreshToken, err := h.TokenM.GenerateNew(userID, clientIP)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	err = h.Repo.Insert(userID, clientIP, refreshToken, h.TokenM.HashAccessToken(accessToken))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to insert token to DataBase")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *Handler) HandleUpdateTokens(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusBadRequest, "Invalid http method")
		return
	}

	log.Printf("Handle update is active. CLIENT: %s", r.RemoteAddr)

	var req struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	clientIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid client IP address")
		return
	}

	parsedUserID, err := h.TokenM.ParseAccess(req.AccessToken, "user_id")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid access token")
		return
	}

	parsedClientIP, err := h.TokenM.ParseAccess(req.AccessToken, "ip")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid access token")
		return
	}

	if clientIP != parsedClientIP {
		if err := mail.SendToIP(clientIP, "Warning. Someone has access to your account."); err != nil {
			log.Printf("Failed to send email to %s", clientIP)
		}
	}

	ok, err := h.Repo.Check(parsedUserID, req.RefreshToken, h.TokenM.HashAccessToken(req.AccessToken))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to check refresh token")
		return
	}
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Invalid refresh token")
		return
	}

	if err := h.Repo.MarkUsed(parsedUserID, req.RefreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update refresh token")
		return
	}

	accessToken, refreshToken, err := h.TokenM.GenerateNew(parsedUserID, clientIP)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update tokens")
		return
	}

	if err := h.Repo.Insert(parsedUserID, clientIP, refreshToken, h.TokenM.HashAccessToken(accessToken)); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to insert token to DataBase")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
