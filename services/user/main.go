package user

import (
	"context"

	"DouyinMerchant/pkg/middleware/auth"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
)

func main() {
	h := server.Default()

	// Initialize Redis for token blacklist
	auth.InitRedis()

	// Apply middleware globally
	h.Use(auth.AuthMiddleware())

	// Public routes
	public := h.Group("/api")
	{
		public.POST("/user/register", registerHandler)
		public.POST("/user/login", loginHandler)
	}

	// Protected routes
	protected := h.Group("/api")
	{
		protected.GET("/user/profile", getUserProfile)
		protected.POST("/user/logout", logoutHandler)
	}

	h.Spin()
}

func getUserProfile(ctx context.Context, ctx2 *app.RequestContext) {

}

func registerHandler(ctx context.Context, ctx2 *app.RequestContext) {

}

func loginHandler(ctx context.Context, c *app.RequestContext) {
	// After successful login
	userId := 1
	token, err := auth.GenerateToken(int32(userId))
	if err != nil {
		// Handle error
		return
	}

	c.JSON(200, map[string]interface{}{
		"token": token,
	})
}

func logoutHandler(ctx context.Context, c *app.RequestContext) {
	token := c.Request.Header.Get("Authorization")
	err := auth.BlacklistToken(token)
	if err != nil {
		// Handle error
		return
	}

	c.JSON(200, map[string]interface{}{
		"message": "Successfully logged out",
	})
}
