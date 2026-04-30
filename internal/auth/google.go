package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// GoogleTokenInfo adalah response dari Google tokeninfo API
type GoogleTokenInfo struct {
	Sub           string `json:"sub"`   // Google User ID
	Aud           string `json:"aud"`   // Client ID yang digunakan saat generate token
	Email         string `json:"email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	EmailVerified string `json:"email_verified"`
	Error         string `json:"error"`
}

// VerifyGoogleToken memverifikasi Google ID token via Google tokeninfo endpoint.
// Jika GOOGLE_CLIENT_ID diset di env, validasi audience (aud) dilakukan
// untuk memastikan token berasal dari aplikasi Android yang benar.
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

	// Validasi audience: pastikan token dibuat oleh aplikasi Android kita
	// Jika GOOGLE_CLIENT_ID tidak diset di env, skip validasi ini (mode dev)
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID != "" && info.Aud != clientID {
		return nil, errors.New("token Google tidak valid: audience tidak sesuai")
	}

	return &info, nil
}
