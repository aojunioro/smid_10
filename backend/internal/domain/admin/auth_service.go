package admin

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserInactive       = errors.New("user is inactive")
)

// AuthService define operações de autenticação.
type AuthService struct {
	userRepo    UserRepository
	jwtSecret   string
	jwtExpiry   time.Duration
}

// NewAuthService cria uma nova instância de AuthService.
func NewAuthService(userRepo UserRepository, jwtSecret string, jwtExpiry time.Duration) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

// LoginRequest representa uma requisição de login.
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginResponse representa a resposta de um login bem-sucedido.
type LoginResponse struct {
	Token     string      `json:"token"`
	ExpiresAt time.Time   `json:"expires_at"`
	User      SystemUser  `json:"user"`
}

// Login autentica um usuário e retorna um token JWT.
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	// Buscar usuário pelo login
	user, err := s.userRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Verificar se o usuário está ativo
	if !user.Active {
		return nil, ErrUserInactive
	}

	// Verificar senha
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Gerar token JWT
	token, expiresAt, err := s.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Remover hash da senha da resposta
	user.PasswordHash = ""

	return &LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *user,
	}, nil
}

// Claims representa os claims do token JWT.
type Claims struct {
	UserID   int64  `json:"user_id"`
	Login    string `json:"login"`
	Name     string `json:"name"`
	UnitID   *int64 `json:"unit_id"`
	jwt.RegisteredClaims
}

// generateToken gera um token JWT para o usuário.
func (s *AuthService) generateToken(user *SystemUser) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(s.jwtExpiry)

	claims := Claims{
		UserID: user.ID,
		Login:  user.Login,
		Name:   user.Name,
		UnitID: user.SystemUnitID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Issuer:    "smid10",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken valida um token JWT e retorna os claims.
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validar método de assinatura
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidCredentials
}

// HashPassword gera um hash bcrypt da senha.
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// ComparePassword compara uma senha em texto plano com um hash bcrypt.
func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// ConstantTimeCompare compara duas strings em tempo constante para evitar timing attacks.
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
