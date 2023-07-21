package helpers

type SignUpInput struct {
	Email       string `json:"email" validate:"required"`
	Password    string `json:"password" validate:"required,min=8"`
	LastName    string `json:"last_name" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	UserName    string `json:"user_name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}

type VerifyEmailInput struct {
	Email string `json:"email" validate:"required"`
	OTP   int `json:"otp" validate:"required"`
}

type ResendOTPInput struct {
	Email string `json:"email" validate:"required"`
}

type LoginInput struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type IError struct {
	Field string
	Tag   string
	Value string
}
