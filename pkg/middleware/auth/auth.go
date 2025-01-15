package auth

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type AuthMiddleware struct {
	jwtSecret []byte
	Rdb       *redis.Client
	whitelist map[string]bool
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func NewAuthMiddleware(secret string, rdb *redis.Client) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(secret),
		Rdb:       rdb,
		whitelist: map[string]bool{
			"/api/user/register": true,
			"/api/user/login":    true,
		},
	}
}

func (m *AuthMiddleware) GenerateToken(username string) (string, error) {
	now := time.Now()
	claims := Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(24 * time.Hour).Unix(),
			IssuedAt:  now.Unix(),
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.jwtSecret)
}

func (m *AuthMiddleware) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return m.jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// Check if token is blacklisted
		exists, err := m.Rdb.Exists(context.Background(), "blacklist:"+tokenString).Result()
		if err != nil || exists > 0 {
			return nil, jwt.ErrInvalidKey
		}
		return claims, nil
	}

	return nil, jwt.ErrInvalidKey
}

func (m *AuthMiddleware) BlacklistToken(tokenString string, expiration time.Duration) error {
	return m.Rdb.Set(context.Background(), "blacklist:"+tokenString, true, expiration).Err()
}

// Hertz middleware function
func (m *AuthMiddleware) AuthRequired() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())

		// Check whitelist
		if m.whitelist[path] {
			c.Next(ctx)
			return
		}

		// Get token from header
		authHeader := string(c.Request.Header.Peek("Authorization"))
		if authHeader == "" {
			c.AbortWithStatus(401)
			return
		}

		// Parse token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatus(401)
			return
		}

		claims, err := m.ValidateToken(parts[1])
		if err != nil {
			c.AbortWithStatus(401)
			return
		}

		// Set username in context
		c.Set("username", claims.Username)
		c.Next(ctx)
	}
}

// HashPassword hashes the given plaintext password using bcrypt.
// It returns the hashed password as a string and any error encountered.
func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash of the password with the default cost.
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	// Convert the hashed bytes to a string and return.
	return string(hashedBytes), nil
}

// CheckPasswordHash compares a plaintext password with a bcrypt hashed password.
// It returns true if they match, and false otherwise.
func CheckPasswordHash(password, hashedPassword string) bool {
	// Compare the hashed password with the plaintext password.
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	// If there's no error, the passwords match.
	return err == nil
}
