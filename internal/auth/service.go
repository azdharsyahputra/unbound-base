package auth

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	DB        *gorm.DB
	JWTSecret []byte
}

func NewAuthService(db *gorm.DB) *AuthService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	return &AuthService{DB: db, JWTSecret: []byte(secret)}
}

type RegisterReq struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResp struct {
	Token string `json:"token"`
}

func (s *AuthService) Register(input RegisterReq) (*User, error) {
	if input.Username == "" || input.Email == "" || input.Password == "" {
		return nil, errors.New("username, email, and password are required")
	}

	var count int64
	s.DB.Model(&User{}).Where("email = ? OR username = ?", input.Email, input.Username).Count(&count)
	if count > 0 {
		return nil, errors.New("email/username already used")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hash),
	}

	if err := s.DB.Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

func (s *AuthService) Login(input LoginReq) (*TokenResp, error) {
	var u User
	if err := s.DB.Where("email = ?", input.Email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatUint(uint64(u.ID), 10), // sub = userID
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "unbound",
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := tok.SignedString(s.JWTSecret)
	if err != nil {
		return nil, err
	}

	return &TokenResp{Token: ss}, nil
}

func (s *AuthService) ParseToken(tokenStr string) (uint, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.JWTSecret, nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := tok.Claims.(*jwt.RegisteredClaims); ok && tok.Valid {
		id64, err := strconv.ParseUint(claims.Subject, 10, 64)
		if err != nil {
			return 0, err
		}
		return uint(id64), nil
	}
	return 0, errors.New("invalid token")
}
