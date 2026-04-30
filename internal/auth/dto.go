package auth

// GoogleLoginRequest login user Android via Google ID Token
type GoogleLoginRequest struct {
	IDToken string `json:"id_token" binding:"required"`
}

// AdminLoginRequest login admin dengan username + password
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse response setelah login berhasil
type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}
