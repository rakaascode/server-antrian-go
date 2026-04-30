package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// GoogleTokenInfo adalah response dari Google tokeninfo API
type GoogleTokenInfo struct {
	Sub           string `json:"sub"`   // Google User ID
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified string `json:"email_verified"`
	Error         string `json:"error"`
}

// VerifyGoogleToken memverifikasi Google ID token via Google tokeninfo endpoint
func VerifyGoogleToken(idToken string) (*GoogleTokenInfo, error) {
	url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("gagal verifikasi token Google: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info GoogleTokenInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	if info.Error != "" {
		return nil, errors.New("Google token tidak valid: " + info.Error)
	}

	if info.EmailVerified != "true" {
		return nil, errors.New("email Google belum diverifikasi")
	}

	return &info, nil
}
