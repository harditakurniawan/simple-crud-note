package utils

import (
	"crypto/rsa"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService interface {
	GenerateToken(userID string, userEmail string) (string, error)
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserInfoFromToken(token *jwt.Token) (*UserInfo, error)
}

type jwtService struct {
	privateKey    *rsa.PrivateKey
	publicKey     *rsa.PublicKey
	tokenDuration time.Duration
}

type UserInfo struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func NewJWTService(privateKeyPath, publicKeyPath string, tokenDuration time.Duration) (JWTService, error) {
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		return nil, err
	}

	return &jwtService{
		privateKey:    privateKey,
		publicKey:     publicKey,
		tokenDuration: tokenDuration,
	}, nil
}

func (j *jwtService) GenerateToken(userID string, userEmail string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   userEmail,
		"exp":     time.Now().Add(j.tokenDuration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(j.privateKey)
}

func (j *jwtService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return j.publicKey, nil
	})
}

func (j *jwtService) GetUserInfoFromToken(token *jwt.Token) (*UserInfo, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok1 := claims["user_id"].(string)
		email, ok2 := claims["email"].(string)

		if !ok1 || !ok2 {
			return nil, jwt.ErrInvalidKey
		}

		return &UserInfo{
			UserID: userID,
			Email:  email,
		}, nil
	}
	return nil, jwt.ErrInvalidKey
}
