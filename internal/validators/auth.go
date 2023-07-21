package validators

import (
	"github.com/gofiber/fiber/v2"
	"github.com/showbaba/go-auth-service/internal/helpers"
	"github.com/go-playground/validator"
)

var Validator = validator.New()

func ValidateSignup(c *fiber.Ctx) error {
	body := new(helpers.SignUpInput)
	c.BodyParser(&body)

	err := Validator.Struct(body)
	if err != nil {
		return helpers.SchemaError(c, err)
	}
	return c.Next()
}

func ValidateLogin(c *fiber.Ctx) error {
	body := new(helpers.LoginInput)
	c.BodyParser(&body)

	err := Validator.Struct(body)
	if err != nil {
		return helpers.SchemaError(c, err)
	}
	return c.Next()
}

func ValidateVerifyEmail(c *fiber.Ctx) error {
	body := new(helpers.VerifyEmailInput)
	c.BodyParser(&body)

	err := Validator.Struct(body)
	if err != nil {
		return helpers.SchemaError(c, err)
	}
	return c.Next()
}
