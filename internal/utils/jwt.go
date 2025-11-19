package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vmaurya-21/Calance-Workflow/internal/config"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID      uuid.UUID `json:"user_id"`
	GitHubID    int64     `json:"github_id"`
	Username    string    `json:"username"`
	Email       string    `json:"email"`
	GitHubToken string    `json:"github_token"`
	jwt.RegisteredClaims
}

// GenerateToken generates a new JWT token for a user
func GenerateToken(userID uuid.UUID, githubID int64, username, email, githubToken string) (string, error) {
	cfg := config.AppConfig
	if cfg == nil {
		return "", errors.New("config not initialized")
	}

	expirationTime := time.Now().Add(time.Duration(cfg.JWT.ExpirationHours) * time.Hour)

	claims := &JWTClaims{
		UserID:      userID,
		GitHubID:    githubID,
		Username:    username,
		Email:       email,
		GitHubToken: githubToken,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "calance-workflow",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	cfg := config.AppConfig
	if cfg == nil {
		return nil, errors.New("config not initialized")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// ExtractTokenFromHeader extracts the token from Authorization header
func ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	// Expected format: "Bearer <token>"
	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	return authHeader[len(bearerPrefix):], nil
}
