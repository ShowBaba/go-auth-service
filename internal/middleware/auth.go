package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/showbaba/go-auth-service/utils"
)

type AuthMiddleware struct{}

func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func (m *AuthMiddleware) ValidateAuthHeaderToken(c *fiber.Ctx) error {
	tokenInHeader := c.Get("Authorization")
	if tokenInHeader == "" {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"status":  false,
				"message": "missing auth token in header",
			})
	}
	token := strings.Split(tokenInHeader, " ")[1]
	if token == "" {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"status":  false,
				"message": "missing auth token in header",
			})
	}
	claim, err := utils.ValidateAuthToken(token, utils.GetConfig().JWTSecretKey)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{
				"status":  false,
				"message": fmt.Sprintf("error validating auth token token: %v", err),
			})
	}

	// TODO: fetch user and confirm if email has been validated
	c.Set("email", claim.Email)
	c.Set("id", strconv.FormatUint(uint64(claim.ID), 10))
	return c.Next()
}
