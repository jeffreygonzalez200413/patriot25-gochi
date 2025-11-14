package api

import (
	"context"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/juhun32/patriot25-gochi/go/google"
	"github.com/juhun32/patriot25-gochi/go/models"
	"github.com/juhun32/patriot25-gochi/go/repo"
)

type AuthHandler struct {
	Google      *google.GoogleOAuth
	UserRepo    *repo.UserRepo
	JWTSecret   string
	FrontendURL string
}

type Claims struct {
	UserID string `json:"userId"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthHandler(google *google.GoogleOAuth, userRepo *repo.UserRepo, jwtSecret, frontendURL string) *AuthHandler {
	return &AuthHandler{
		Google:      google,
		UserRepo:    userRepo,
		JWTSecret:   jwtSecret,
		FrontendURL: frontendURL,
	}
}

// In real apps, you should generate & verify state to prevent CSRF.
// For hackathon, we keep it simple.
func (h *AuthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := "random-state" // TODO: generate & store in cookie/session
	url := h.Google.AuthCodeURL(state)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *AuthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	if errMsg := r.URL.Query().Get("error"); errMsg != "" {
		http.Error(w, "Google error: "+errMsg, http.StatusBadRequest)
		return
	}

	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "missing code", http.StatusBadRequest)
		return
	}

	userInfo, _, err := h.Google.GetUserInfo(ctx, code)
	if err != nil {
		http.Error(w, "failed to get userinfo: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user := &models.User{
		UserID:  userInfo.Sub,
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Picture: userInfo.Picture,
	}

	if err := h.UserRepo.UpsertUser(ctx, user); err != nil {
		http.Error(w, "failed to save user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate JWT
	token, err := GenerateJWT(h.JWTSecret, user.UserID, user.Email, 7*24*time.Hour)
	if err != nil {
		http.Error(w, "failed to generate token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set cookie (adjust Secure/SameSite in production)
	http.SetCookie(w, &http.Cookie{
		Name:     "ppet_token",
		Value:    token,
		Path:     "/",
		HttpOnly: false,
		Secure:   false, // true in HTTPS/prod
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int((7 * 24 * time.Hour).Seconds()),
	})

	// Redirect back to your frontend app
	redirect := h.FrontendURL
	if redirect == "" {
		redirect = "http://localhost:3000/app"
	}
	http.Redirect(w, r, redirect, http.StatusFound)
}

// JWT generation and validation functions
func GenerateJWT(secret, userID, email string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseAndValidateJWT(secret, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}
