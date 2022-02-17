package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/kazakovichna/todoListPrjct"
	"github.com/kazakovichna/todoListPrjct/pkg/repository"
	"math/rand"
	"time"
)

const (
	salt = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL = 12 * time.Hour
)


func PlusSomething(r int) int {
	ans := r + r
	return ans
}


type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todoListPrjct.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)

	return s.repo.CreateUser(user)
}

func (s *AuthService) CreateSession(username, password string) (TokenResponse, error) {
	var (
		res TokenResponse
		err error
	)

	res.AccessToken, err = s.GenerateToken(username, password)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.GenerateRefreshToken()
	if err != nil {
		return res, nil
	}

	expiresAt := time.Now().Add(3 *time.Minute).Unix()

	err = s.repo.SetSessions(username, res.RefreshToken, int(expiresAt))

	return res, nil
}

type TokenResponse struct {
	AccessToken string
	RefreshToken string
}

func (s *AuthService) SignInService(username, password string) (TokenResponse, error) {
	return s.CreateSession(username, password)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 *time.Minute).Unix(),
			IssuedAt: time.Now().Unix(),
		},
	user.Id,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) GenerateRefreshToken() (string, error) {
	b := make([]byte, 32)

	m := rand.NewSource(time.Now().Unix())
	r := rand.New(m)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (s *AuthService) RefreshTokenServices(refreshToken string) (TokenResponse, error) {
	userInfo, err := s.repo.GetUserByRefreshToken(refreshToken)
	if err != nil && userInfo.RefreshToken != refreshToken {
		return TokenResponse{
			"refresh token are not valid",
			"get a fuck out here or sign in again",
		}, err
	}

	return s.CreateSessionByRefreshToken(userInfo)
}

func (s *AuthService) CreateSessionByRefreshToken(user todoListPrjct.UserRefreshToken) (TokenResponse, error) {
	var (
		res TokenResponse
		err error
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 *time.Minute).Unix(),
			IssuedAt: time.Now().Unix(),
		},
		user.Id,
	})

	res.AccessToken, err = token.SignedString([]byte(signingKey))
	if err != nil {
		return TokenResponse{}, err
	}

	res.RefreshToken, err = s.GenerateRefreshToken()
	if err != nil {
		return TokenResponse{}, err
	}

	expiresAt := time.Now().Add(3 *time.Minute).Unix()

	err = s.repo.SetSessions(user.Username, res.RefreshToken, int(expiresAt))

	return res, nil
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
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


