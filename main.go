package main

import (
	"DouyinMerchant/pkg/middleware/auth"
	"context"
	"github.com/go-redis/redis/v8"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

var authMiddleware *auth.AuthMiddleware

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	authMiddleware = auth.NewAuthMiddleware("your-secret-key", rdb)

	h := server.Default()
	// Apply middleware globally
	h.Use(authMiddleware.AuthRequired())

	// Routes
	h.POST("/api/user/register", registerHandler)
	h.POST("/api/user/login", loginHandler)

	h.Spin()
}

func loginHandler(ctx context.Context, ctx2 *app.RequestContext) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx2.Bind(&req); err != nil {
		ctx2.JSON(400, map[string]string{"error": "Invalid request"})
		return
	}

	if req.Username == "" || req.Password == "" {
		ctx2.JSON(400, map[string]string{"error": "Username and password are required"})
		return
	}
	// Get user from Redis
	hashedPassword, err := authMiddleware.Rdb.Get(ctx, req.Username).Result()
	if err == redis.Nil {
		ctx2.JSON(400, map[string]string{"error": "Invalid username or password"})
		return
	} else if err != nil {
		ctx2.JSON(500, map[string]string{"error": "Failed to retrieve user"})
		return
	}

	// Compare passwords
	if !auth.CheckPasswordHash(req.Password, hashedPassword) {
		ctx2.JSON(400, map[string]string{"error": "Invalid username or password"})
		return
	}

	// Generate JWT token
	authMiddleware := auth.NewAuthMiddleware("your-secret-key", authMiddleware.Rdb) // Ensure you have access to AuthMiddleware
	token, err := authMiddleware.GenerateToken(req.Username)
	if err != nil {
		ctx2.JSON(500, map[string]string{"error": "Failed to generate token"})
		return
	}

	ctx2.JSON(200, map[string]string{"message": "Login successful", "token": token})
}

func registerHandler(ctx context.Context, ctx2 *app.RequestContext) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx2.Bind(&req); err != nil {
		ctx2.JSON(400, map[string]string{"error": "Invalid request"})
		return
	}

	if req.Username == "" || req.Password == "" {
		ctx2.JSON(400, map[string]string{"error": "Username and password are required"})
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		ctx2.JSON(500, map[string]string{"error": "Failed to hash password"})
		return
	}

	// Check if username exists in Redis
	if authMiddleware.Rdb.Get(ctx, req.Username).Err() != redis.Nil {
		ctx2.JSON(400, map[string]string{"error": "Username already exists"})
		return
	}

	// Store user in Redis
	if err := authMiddleware.Rdb.Set(ctx, req.Username, hashedPassword, 0).Err(); err != nil {
		ctx2.JSON(500, map[string]string{"error": "Failed to save user"})
		return
	}

	ctx2.JSON(200, map[string]string{"message": "User registered successfully"})
}
