package handler

import (
	"net/http"
	"fmt"
)

func HandleGenerateTokens(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "Missing user_id", http.StatusBadRequest)
		return
	}
	fmt.Println(userID)
	// accessToken, refreshToken, err := token.GenerateNew(userID)
}

func HandleUpdateTokens(w http.ResponseWriter, r *http.Request) {
	return
}
