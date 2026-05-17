package route

import (
	"github.com/giakiet05/lkforum/internal/controller"
	"github.com/giakiet05/lkforum/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
)

// RegisterAuthRoutes registers all authentication-related routes.
func RegisterAuthRoutes(rg *gin.RouterGroup, authCtrl *controller.AuthController, userCtrl *controller.UserController, redisClient *redis.Client) {
	auth := rg.Group("/auth")

	auth.POST("/refresh", authCtrl.RefreshToken)
	auth.POST("/logout", authCtrl.Logout)
	auth.POST("/check-username", userCtrl.CheckUsername) // Public endpoint for username availability check

	// Local Authentication - New Flow (Verify Email First)
	local := auth.Group("/local")
	{
		local.POST("/send-verification", middleware.RateLimit(redisClient, "auth_send_verification", 5, 10*time.Minute), authCtrl.SendEmailVerification)
		local.POST("/verify-email", middleware.RateLimit(redisClient, "auth_verify_email", 10, 10*time.Minute), authCtrl.VerifyEmailCode)
		local.POST("/complete-registration", authCtrl.CompleteRegistration)
		local.POST("/resend-otp", middleware.RateLimit(redisClient, "auth_resend_otp", 5, 10*time.Minute), authCtrl.ResendOTP)
		local.POST("/login", middleware.RateLimit(redisClient, "auth_login", 10, 5*time.Minute), authCtrl.Login)

		// Forgot Password Flow (only for local auth)
		local.POST("/forgot-password", middleware.RateLimit(redisClient, "auth_forgot_password", 5, 10*time.Minute), authCtrl.ForgotPassword)
		local.POST("/verify-reset-otp", middleware.RateLimit(redisClient, "auth_verify_reset_otp", 10, 10*time.Minute), authCtrl.VerifyResetPasswordOTP)
		local.POST("/reset-password", authCtrl.ResetPassword)
	}

	// Google OAuth2
	google := auth.Group("/google")
	{
		google.GET("/login", authCtrl.GoogleLogin)
		google.GET("/callback", authCtrl.GoogleCallback)
		google.POST("/complete-setup", authCtrl.CompleteGoogleSetup)
	}
}
