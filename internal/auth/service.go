package auth

import (
	"errors"

	"github.com/rakaascode/server-antrian-go.git/internal/user"
	"github.com/rakaascode/server-antrian-go.git/pkg/utils"
)

type AuthService interface {
	// GoogleLogin — khusus user Android
	GoogleLogin(req GoogleLoginRequest) (AuthResponse, error)
	// AdminLogin — khusus admin cabang (username + password)
	AdminLogin(req AdminLoginRequest) (AuthResponse, error)
}

type authService struct {
	userRepo user.UserRepository
}

func NewAuthService(repo user.UserRepository) AuthService {
	return &authService{repo}
}

// GoogleLogin — login/register user Android via Google ID Token
func (s *authService) GoogleLogin(req GoogleLoginRequest) (AuthResponse, error) {
	info, err := VerifyGoogleToken(req.IDToken)
	if err != nil {
		return AuthResponse{}, err
	}

	// Cari user berdasarkan Google ID
	u, err := s.userRepo.FindByGoogleID(info.Sub)
	if err != nil {
		// Belum ada, coba cari by email
		u, err = s.userRepo.FindByEmail(info.Email)
		if err != nil {
			// Buat akun baru otomatis
			googleID := info.Sub
			emailStr := info.Email
			newUser := user.User{
				Name:      info.Name,
				Email:     &emailStr,
				GoogleID:  &googleID,
				AvatarURL: info.Picture,
				Role:      "user",
			}
			u, err = s.userRepo.Create(newUser)
			if err != nil {
				return AuthResponse{}, errors.New("gagal membuat akun: " + err.Error())
			}
		} else {
			// Update GoogleID dan AvatarURL jika belum ada atau berubah
			updated := false
			if u.GoogleID == nil {
				googleID := info.Sub
				u.GoogleID = &googleID
				updated = true
			}
			if u.AvatarURL != info.Picture && info.Picture != "" {
				u.AvatarURL = info.Picture
				updated = true
			}
			if updated {
				u, _ = s.userRepo.Update(u)
			}
		}
	} else {
		// User ditemukan by GoogleID — update avatar jika berubah
		if u.AvatarURL != info.Picture && info.Picture != "" {
			u.AvatarURL = info.Picture
			u, _ = s.userRepo.Update(u)
		}
	}

	emailStr := ""
	if u.Email != nil {
		emailStr = *u.Email
	}
	token, err := GenerateToken(u.ID, emailStr, "", u.Role, nil)
	if err != nil {
		return AuthResponse{}, err
	}
	u.Password = ""
	return AuthResponse{Token: token, User: u}, nil
}

// AdminLogin — login admin cabang dengan username + password
func (s *authService) AdminLogin(req AdminLoginRequest) (AuthResponse, error) {
	u, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return AuthResponse{}, errors.New("username atau password salah")
	}
	if u.Role != "admin" {
		return AuthResponse{}, errors.New("akun ini bukan admin")
	}
	if !utils.CheckPasswordHash(req.Password, u.Password) {
		return AuthResponse{}, errors.New("username atau password salah")
	}
	uname := ""
	if u.Username != nil {
		uname = *u.Username
	}
	token, err := GenerateToken(u.ID, "", uname, u.Role, u.CabangID)
	if err != nil {
		return AuthResponse{}, err
	}
	u.Password = ""
	return AuthResponse{Token: token, User: u}, nil
}
