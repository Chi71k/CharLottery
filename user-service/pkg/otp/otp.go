package otp

import (
	"crypto/rand"
	"fmt"
	"time"
)

// Генерация случайного OTP
func GenerateOTP() string {
	var otp [6]byte
	_, err := rand.Read(otp[:])
	if err != nil {
		panic("failed to generate OTP")
	}
	return fmt.Sprintf("%06d", otp)
}

// Генерация уникального OTP с использованием текущего времени
func GenerateUniqueOTP() string {
	// Используем текущую метку времени для генерации уникального OTP
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%x", timestamp)
}
