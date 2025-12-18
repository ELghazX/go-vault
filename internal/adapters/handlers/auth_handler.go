package handlers

import (
	"net/http"
	"time"

	"github.com/elghazx/go-vault/internal/core/services"
	"github.com/elghazx/go-vault/templates"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	token, err := h.authService.Login(c.Request().Context(), username, password)
	if err != nil {
		return templates.ErrorMessage("Invalid credentials").Render(c.Request().Context(), c.Response().Writer)
	}

	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
	}

	c.SetCookie(cookie)
	return templates.SuccessMessage("Login successful").Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) Register(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Invalid request")
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	if err := h.authService.Register(c.Request().Context(), username, password); err != nil {
		return templates.ErrorMessage("Registration failed").Render(c.Request().Context(), c.Response().Writer)
	}

	return templates.SuccessMessage("Registration successful").Render(c.Request().Context(), c.Response().Writer)
}

func (h *AuthHandler) CheckAuth(c echo.Context) error {
	cookie, err := c.Cookie("auth_token")
	if err != nil {
		return c.String(http.StatusUnauthorized, "Not authenticated")
	}

	_, err = h.authService.ValidateToken(cookie.Value)
	if err != nil {
		return c.String(http.StatusUnauthorized, "Invalid token")
	}

	return c.String(http.StatusOK, "authenticated")
}

func (h *AuthHandler) Logout(c echo.Context) error {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
	}

	c.SetCookie(cookie)
	return templates.SuccessMessage("Logged out successfully").Render(c.Request().Context(), c.Response().Writer)
}
