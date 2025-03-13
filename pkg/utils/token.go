package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
type JWTMaker struct {
	secretKey string
}

// Payload contains the payload data of the JWT
type Payload struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// NewJWTMaker creates a new JWTMaker
func NewJWTMaker(secretKey string) (*JWTMaker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a user
func (maker *JWTMaker) CreateToken(userID uuid.UUID, username string, role string, duration time.Duration) (string, error) {
	payload := &Payload{
		ID:        uuid.New(),
		UserID:    userID,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         payload.ID.String(),
		"user_id":    payload.UserID.String(),
		"username":   payload.Username,
		"role":       payload.Role,
		"issued_at":  payload.IssuedAt.Unix(),
		"expired_at": payload.ExpiredAt.Unix(),
	})

	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token signing method")
		}
		return []byte(maker.secretKey), nil
	}

	jwtToken, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok || !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	// Parse token claims
	idStr, _ := claims["id"].(string)
	userIDStr, _ := claims["user_id"].(string)
	username, _ := claims["username"].(string)
	role, _ := claims["role"].(string)
	issuedAt, _ := claims["issued_at"].(float64)
	expiredAt, _ := claims["expired_at"].(float64)

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, errors.New("invalid token ID")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID")
	}

	payload := &Payload{
		ID:        id,
		UserID:    userID,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Unix(int64(issuedAt), 0),
		ExpiredAt: time.Unix(int64(expiredAt), 0),
	}

	// Check if the token has expired
	if time.Now().After(payload.ExpiredAt) {
		return nil, errors.New("token has expired")
	}

	return payload, nil
}
