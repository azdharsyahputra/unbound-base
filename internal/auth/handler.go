package auth

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterRoutes(app *fiber.App, db *gorm.DB, svc *AuthService) {
	r := app.Group("/auth")

	r.Post("/register", func(c *fiber.Ctx) error {
		return RegisterHandler(c, svc)
	})
	r.Post("/login", func(c *fiber.Ctx) error {
		return LoginHandler(c, svc)
	})
	r.Post("/refresh", func(c *fiber.Ctx) error {
		return RefreshHandler(c, svc)
	})
	r.Post("/logout", func(c *fiber.Ctx) error {
		return LogoutHandler(c, svc)
	})
}

// ================================================================
// HANDLERS
// ================================================================

// RegisterHandler godoc
// @Summary Register new user
// @Description Membuat akun baru di Unbound
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body RegisterReq true "User data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/register [post]
func RegisterHandler(c *fiber.Ctx, svc *AuthService) error {
	var req RegisterReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}
	u, err := svc.Register(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
	})
}

// LoginHandler godoc
// @Summary Login user
// @Description Login dengan email dan password untuk mendapatkan token
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body LoginReq true "User credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login [post]
func LoginHandler(c *fiber.Ctx, svc *AuthService) error {
	var req LoginReq
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	tok, err := svc.Login(req)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":        true,
		"access_token":   tok.AccessToken,
		"refresh_token":  tok.RefreshToken,
		"token_type":     "Bearer",
		"expires_in_sec": 86400,
	})
}

// RefreshHandler godoc
// @Summary Refresh token
// @Description Mengambil access_token baru menggunakan refresh_token
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body map[string]string true "Refresh token"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/refresh [post]
func RefreshHandler(c *fiber.Ctx, svc *AuthService) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid body")
	}

	tok, err := svc.RefreshAccess(body.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}

	return c.JSON(fiber.Map{
		"success":        true,
		"access_token":   tok.AccessToken,
		"refresh_token":  tok.RefreshToken,
		"token_type":     "Bearer",
		"expires_in_sec": 86400,
	})
}

// LogoutHandler godoc
// @Summary Logout user
// @Description Menghapus refresh_token dari database (invalidate session)
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body map[string]string true "Refresh token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /auth/logout [post]
func LogoutHandler(c *fiber.Ctx, svc *AuthService) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&body); err != nil || body.RefreshToken == "" {
		return fiber.NewError(fiber.StatusBadRequest, "refresh_token required")
	}

	if err := svc.DB.Where("token = ?", body.RefreshToken).Delete(&RefreshToken{}).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to logout")
	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "Logged out successfully",
	})
}
