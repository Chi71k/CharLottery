package model

// User - модель пользователя
type User struct {
    ID       string `json:"id" bson:"_id,omitempty"`
    Username string `json:"username" bson:"username"`
    Email    string `json:"email" bson:"email"`
    Password string `json:"password" bson:"password"`
}

// RegisterRequest - структура запроса для регистрации
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// VerifyOTPRequest - структура запроса для верификации OTP
type VerifyOTPRequest struct {
	UserID string `json:"user_id"`
	Otp    string `json:"otp"`
}
