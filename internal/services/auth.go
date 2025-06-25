package services

import (
	"merch-store/internal/repository"
	"merch-store/models"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) FindUser(username string) (int, error) {
    return s.repo.GetUserID(username)
}

func (s *AuthService) SignIn(username string, password string) (int, error) {
	hashPassword, err := s.repo.GetUserPassword(username)
	if err != nil {
		return 0, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password)); err != nil {
		return 0, err
	}

	return s.repo.GetUserID(username)
}

func (s *AuthService) CreateUser(username string, password string) (int, error) {
	password, err := hashPassword(password)
	if err != nil {
		return 0, err 
	}
	user := models.AuthRequest {
		Username: username,
		Password: password,
	}

	id, err := s.repo.CreateUser(user)
	if err != nil {
		log.Println("Не удалось добавить пользователя в базу")
		return 0, err
	}

	user_wallet := models.Wallet {
		User_id: id,
	}
	err = s.repo.AddCoins(user_wallet)
	if err != nil {
		log.Println("Ошибка при начислении коинов")
		return 0, err
	}

	return id, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}


func (s *AuthService) GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userID,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}
