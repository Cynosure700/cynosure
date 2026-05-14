package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"nano_cc/internal/config"
	"nano_cc/internal/web/storage"
)

type Service struct {
	Store *storage.Store
	Cfg   config.AppConfig
}

type Claims struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	jwt.RegisteredClaims
}

type ContextKey string

const UserContextKey ContextKey = "web_user"

func NewService(store *storage.Store, cfg config.AppConfig) *Service {
	return &Service{Store: store, Cfg: cfg}
}

func (s *Service) Register(ctx context.Context, email, username, password string) (storage.User, string, storage.AuthSession, error) {
	if strings.TrimSpace(email) == "" || strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return storage.User{}, "", storage.AuthSession{}, errors.New("email, username and password are required")
	}
	if _, err := s.Store.GetUserByEmailOrUsername(ctx, email); err == nil {
		return storage.User{}, "", storage.AuthSession{}, errors.New("email already exists")
	}
	if _, err := s.Store.GetUserByEmailOrUsername(ctx, username); err == nil {
		return storage.User{}, "", storage.AuthSession{}, errors.New("username already exists")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return storage.User{}, "", storage.AuthSession{}, err
	}
	user := storage.User{ID: newID("usr"), Email: strings.ToLower(email), Username: username, PasswordHash: string(hash)}
	if err := s.Store.CreateUser(ctx, user); err != nil {
		return storage.User{}, "", storage.AuthSession{}, err
	}
	token, session, err := s.createSessionToken(user.ID)
	if err != nil {
		return storage.User{}, "", storage.AuthSession{}, err
	}
	if err := s.Store.CreateSession(ctx, session); err != nil {
		return storage.User{}, "", storage.AuthSession{}, err
	}
	created, err := s.Store.GetUserByID(ctx, user.ID)
	return created, token, session, err
}

func (s *Service) Login(ctx context.Context, login, password string) (storage.User, string, storage.AuthSession, error) {
	user, err := s.Store.GetUserByEmailOrUsername(ctx, login)
	if err != nil {
		return storage.User{}, "", storage.AuthSession{}, errors.New("invalid credentials")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return storage.User{}, "", storage.AuthSession{}, errors.New("invalid credentials")
	}
	token, session, err := s.createSessionToken(user.ID)
	if err != nil {
		return storage.User{}, "", storage.AuthSession{}, err
	}
	if err := s.Store.CreateSession(ctx, session); err != nil {
		return storage.User{}, "", storage.AuthSession{}, err
	}
	return user, token, session, nil
}

func (s *Service) Logout(ctx context.Context, sessionID string) error {
	return s.Store.RevokeSession(ctx, sessionID)
}

func (s *Service) AuthenticateRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, _, err := s.GetCurrentUser(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (s *Service) GetCurrentUser(r *http.Request) (storage.User, Claims, error) {
	cookie, err := r.Cookie(s.Cfg.CookieName)
	if err != nil || cookie.Value == "" {
		return storage.User{}, Claims{}, errors.New("missing session cookie")
	}
	claims := Claims{}
	token, err := jwt.ParseWithClaims(cookie.Value, &claims, func(token *jwt.Token) (any, error) {
		return []byte(s.Cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		return storage.User{}, Claims{}, errors.New("invalid token")
	}
	session, err := s.Store.GetSession(r.Context(), claims.SessionID)
	if err != nil || session.Status != "active" || session.ExpiresAt.Before(time.Now()) {
		return storage.User{}, Claims{}, errors.New("invalid session")
	}
	user, err := s.Store.GetUserByID(r.Context(), claims.UserID)
	if err != nil {
		return storage.User{}, Claims{}, err
	}
	return user, claims, nil
}

func (s *Service) SetSessionCookie(w http.ResponseWriter, token string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.Cfg.CookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expiresAt,
	})
}

func (s *Service) ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: s.Cfg.CookieName, Path: "/", MaxAge: -1, HttpOnly: true})
}

func (s *Service) createSessionToken(userID string) (string, storage.AuthSession, error) {
	session := storage.AuthSession{
		ID:        newID("sess"),
		UserID:    userID,
		Status:    "active",
		ExpiresAt: time.Now().Add(time.Duration(s.Cfg.SessionTTLMinutes) * time.Minute),
	}
	claims := Claims{
		UserID:    userID,
		SessionID: session.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(session.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "nano-cc-web",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.Cfg.JWTSecret))
	return signed, session, err
}

func UserFromContext(ctx context.Context) (storage.User, bool) {
	user, ok := ctx.Value(UserContextKey).(storage.User)
	return user, ok
}

func newID(prefix string) string {
	raw := make([]byte, 16)
	_, _ = rand.Read(raw)
	mac := hmac.New(sha256.New, raw[:8])
	mac.Write(raw[8:])
	sum := mac.Sum(nil)
	return fmt.Sprintf("%s_%s", prefix, base64.RawURLEncoding.EncodeToString(sum[:12]))
}
