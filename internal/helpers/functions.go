package helpers

import (
	"crypto/rand"
	"math/big"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
)

func SchemaError(c *fiber.Ctx, err error) error {
	var errors []*IError
	for _, err := range err.(validator.ValidationErrors) {
		var el IError
		el.Field = err.Field()
		el.Tag = err.Tag()
		el.Value = err.Param()
		errors = append(errors, &el)
	}
	return c.Status(fiber.StatusBadRequest).JSON(
		&fiber.Map{"errors": errors},
	)
}

func GenerateRandomNumber(min, max int) (int, error) {
	randomNumber, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}
	randomInt := int(randomNumber.Int64()) + min
	return randomInt, nil
}
