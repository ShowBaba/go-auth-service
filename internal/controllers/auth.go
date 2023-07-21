package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rabbitmq/amqp091-go"
	"github.com/showbaba/go-auth-service/internal/helpers"
	"github.com/showbaba/go-auth-service/internal/repository"
	"github.com/showbaba/go-auth-service/models"
	"github.com/showbaba/go-auth-service/utils"
	log "github.com/showbaba/go-auth-service/utils/logging"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var ctx = context.Background()

type AuthController struct {
	database        *gorm.DB
	queueConnection *amqp091.Connection
	userRepository  *repository.UserRepository
	tokenRepository *repository.TokenRepository
}

func NewAuthController(db *gorm.DB, qC *amqp091.Connection,
	tokenRepository *repository.TokenRepository,
	userRepository *repository.UserRepository) *AuthController {
	return &AuthController{
		database:        db,
		queueConnection: qC,
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}
}

func (a *AuthController) Signup(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")

	var input helpers.SignUpInput
	if err := c.BodyParser(&input); err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprintf(`"failed to parse request body: %v`, err)})
	}

	// validate email does not exist
	_, exist, err := a.userRepository.Fetch(models.User{Email: input.Email})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	if exist {
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(
				&fiber.Map{"message": "email already exists"})
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}
	otp, err := helpers.GenerateRandomNumber(10000, 99999)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	fmt.Println("otp value", otp)

	newUser := models.User{
		FirstName:   input.FirstName,
		LastName:    input.LastName,
		Email:       strings.ToLower(input.Email),
		Password:    string(hashedPassword),
		PhoneNumber: input.PhoneNumber,
		Username:    input.UserName,
	}

	user, err := a.userRepository.Create(&newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}
	_, err = a.tokenRepository.Create(&models.Token{UserID: user.ID, Token: otp})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	// send otp to user mail by calling the notification queue
	mail := utils.Mail{
		Sender:  utils.MAIL_USERNAME,
		Subject: "Welcome To Go-Auth-Service!",
		To:      []string{input.Email},
		Body: `<div style="font-family: Helvetica, Arial, sans-serif; min-width: 1000px; overflow: auto; line-height: 2;">
            <div style="margin: 50px auto; width: 70%; padding: 20px 0;">
                <div style="border-bottom: 1px solid #eee;"><a href="google.com" style="font-size: 1.4em; color: #00466a; text-decoration: none; font-weight: 600;">Go-Auth-Service</a></div>
                <p style="font-size: 1.1em;">Hi,</p>
                <p>Hi ` + input.FirstName + `</p>
                <p>Welcome to Go-Auth-Service</p>
                <p>Use the OTP below to verify your account</p>
                <p>` + strconv.Itoa(otp) + `</p>
                <p style="font-size: 0.9em;">
                    Regards,<br />
                    Go-Auth-Service
                </p>
                <hr style="border: none; border-top: 1px solid #eee;" />
            </div>
        </div>`,
	}
	payload, err := json.Marshal(mail)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}
	if err := utils.PublishMessageToQueue(ctx, a.queueConnection, payload, utils.NOTIFICATION_QUEUE); err != nil {
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				&fiber.Map{"message": fmt.Sprint(err)})
		}
	}

	c.Status(201)
	return c.JSON(&fiber.Map{
		"success": true,
		"message": "signup successful",
	})
}

func (a *AuthController) VerifyEmail(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")

	var input helpers.VerifyEmailInput
	if err := c.BodyParser(&input); err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprintf(`"failed to parse request body: %v`, err)})
	}

	user, exist, err := a.userRepository.Fetch(models.User{Email: input.Email})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	if !exist {
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(
				&fiber.Map{"message": "user with email not found"})
		}
	}

	tokenCondition := models.Token{UserID: user.ID, Token: input.OTP}
	token, exist, err := a.tokenRepository.Fetch(tokenCondition)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	if !exist {
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(
				&fiber.Map{"message": "invalid otp"})
		}
	}

	if token == nil {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "invalid token"})
	}

	if token.Token == 0 {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "invalid token"})
	}

	valid := utils.IsTokenValid(*token)

	if !valid {
		return c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "token already expired"})
	}

	err = a.userRepository.Update(user.ID, models.User{IsVerified: utils.BoolPointer(true)})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	jwtToken, err := utils.GenerateToken(utils.GetConfig().JWTSecretKey, input.Email, user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}
	if err = a.tokenRepository.Delete(&tokenCondition); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}
	c.Status(200)
	return c.JSON(&fiber.Map{
		"success": true,
		"message": "email verified successful",
		"token":   jwtToken,
	})
}

func (a *AuthController) ResendOTP(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")

	var input helpers.ResendOTPInput
	if err := c.BodyParser(&input); err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprintf(`"failed to parse request body: %v`, err)})
	}

	user, exist, err := a.userRepository.Fetch(models.User{Email: input.Email})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	if !exist {
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(
				&fiber.Map{"message": "user with email not found"})
		}
	}

	if err = a.tokenRepository.Delete(&models.Token{UserID: user.ID}); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	otp, err := helpers.GenerateRandomNumber(10000, 99999)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	fmt.Println("otp value", otp)
	_, err = a.tokenRepository.Create(&models.Token{UserID: user.ID, Token: otp})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	mail := utils.Mail{
		Sender:  utils.MAIL_USERNAME,
		Subject: "Verify Email OTP!",
		To:      []string{input.Email},
		Body: `<div style="font-family: Helvetica, Arial, sans-serif; min-width: 1000px; overflow: auto; line-height: 2;">
            <div style="margin: 50px auto; width: 70%; padding: 20px 0;">
                <div style="border-bottom: 1px solid #eee;"><a href="google.com" style="font-size: 1.4em; color: #00466a; text-decoration: none; font-weight: 600;">Go-Auth-Service</a></div>
                <p style="font-size: 1.1em;">Hi,</p>
                <p>Hi ` + user.FirstName + `</p>
                <p>Use the OTP below to verify your account</p>
                <p>` + strconv.Itoa(otp) + `</p>
                <p style="font-size: 0.9em;">
                    Regards,<br />
                    Go-Auth-Service
                </p>
                <hr style="border: none; border-top: 1px solid #eee;" />
            </div>
        </div>`,
	}
	payload, err := json.Marshal(mail)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}
	if err := utils.PublishMessageToQueue(ctx, a.queueConnection, payload, utils.NOTIFICATION_QUEUE); err != nil {
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(
				&fiber.Map{"message": fmt.Sprint(err)})
		}
	}

	c.Status(200)
	return c.JSON(&fiber.Map{
		"success": true,
		"message": "otp sent successful",
	})
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")

	var input helpers.LoginInput
	if err := c.BodyParser(&input); err != nil {
		log.Error("failed to parse request body", zap.Error(err))
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprintf(`"failed to parse request body: %v`, err)})
	}

	user, exist, err := a.userRepository.Fetch(models.User{Email: input.Email})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}

	if !exist {
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(
				&fiber.Map{"message": "user with email not found"})
		}
	}

	passwordMatch, err := utils.PasswordMatches(input.Password, user.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}
	if !passwordMatch {
		return c.Status(http.StatusUnauthorized).JSON(
			&fiber.Map{"message": "invalid password"})
	}

	token, err := utils.GenerateToken(utils.GetConfig().JWTSecretKey, input.Email, user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": fmt.Sprint(err)})
	}
	c.Status(200)
	return c.JSON(&fiber.Map{
		"success": true,
		"message": "login successful",
		"token":   token,
	})
}
