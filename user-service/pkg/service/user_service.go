package service

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
	"user-service/pkg/cache"
	"user-service/pkg/mail"
	"user-service/pkg/model"
	"user-service/pkg/natswrap"
	"user-service/pkg/otp"
	"user-service/pkg/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
    Repo  repository.UserRepositoryInterface
    otps  map[string]string
    nats  natswrap.Publisher
    cache cache.CacheInterface
}

func NewUserService(repo repository.UserRepositoryInterface, nats natswrap.Publisher, cache cache.CacheInterface) *UserService {
    return &UserService{
        Repo:  repo,
        otps:  make(map[string]string),
        nats:  nats,
        cache: cache,
    }
}


func (us *UserService) Register(username, password, email string) (string, error) {
	log.Printf("➡️ Register called with: username=%s, email=%s", username, email)

	if us.Repo.UserExists(username) {
		log.Printf("⚠️ Попытка создать уже существующего пользователя: %s", username)
		return "", fmt.Errorf("пользователь уже существует")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("❌ Ошибка при хэшировании пароля: %v", err)
		return "", fmt.Errorf("ошибка при хэшировании пароля: %v", err)
	}

	userID := uuid.New().String()
	user := &model.User{
		ID:       userID,
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	log.Printf("📦 Создание пользователя: %+v", user)

	if err := us.Repo.CreateUser(user); err != nil {
		log.Printf("❌ Ошибка при создании пользователя в БД: %v", err)
		return "", fmt.Errorf("ошибка при создании пользователя: %v", err)
	}

	log.Printf("✅ Пользователь создан: %s (%s)", userID, email)

	otpToken := otp.GenerateUniqueOTP()
	if err := mail.SendOTPEmail(email, otpToken); err != nil {
		log.Printf("❌ Ошибка при отправке OTP: %v", err)
		return "", fmt.Errorf("ошибка при отправке OTP: %v", err)
	}

	log.Printf("📨 OTP отправлен на почту: %s (OTP: %s)", email, otpToken)
	us.otps[userID] = otpToken

	event := struct {
		UserId string `json:"user_id"`
	}{UserId: userID}
	data, _ := json.Marshal(event)
	us.nats.Publish("user.registered", data)

	log.Printf("📣 Event 'user.registered' опубликован для ID: %s", userID)

	return userID, nil
}

func (us *UserService) VerifyOTP(userID, otpCode string) bool {
	log.Printf("➡️ VerifyOTP called: userID=%s, otpCode=%s", userID, otpCode)

	storedOtp, exists := us.otps[userID]
	if !exists || storedOtp != otpCode {
		log.Printf("❌ OTP не совпадает для пользователя %s", userID)
		return false
	}

	delete(us.otps, userID)
	log.Printf("✅ OTP верифицирован для пользователя %s", userID)

	us.nats.Publish("otp.verified", []byte(fmt.Sprintf("OTP verified for UserID: %s", userID)))

	return true
}

func (us *UserService) Login(email, password string) (*model.User, error) {
	log.Printf("➡️ Login attempt: email=%s", email)

	user, err := us.Repo.GetUserByEmail(email)
	if err != nil {
		log.Printf("❌ Пользователь не найден: %s", email)
		return nil, fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Printf("❌ Неверный пароль для пользователя: %s", email)
		return nil, fmt.Errorf("invalid credentials")
	}

	log.Printf("✅ Вход выполнен: %s (%s)", user.ID, email)
	us.nats.Publish("user.logged_in", []byte(fmt.Sprintf("User logged in: %s, Email: %s", user.Username, user.Email)))

	return user, nil
}

func (us *UserService) ForgotPassword(email string) error {
	log.Printf("➡️ ForgotPassword called for: %s", email)

	user, err := us.Repo.GetUserByEmail(email)
	if err != nil {
		log.Printf("❌ Пользователь не найден для сброса пароля: %s", email)
		return fmt.Errorf("user not found")
	}

	otpCode := otp.GenerateUniqueOTP()
	us.otps[user.ID] = otpCode

	if err := mail.SendOTPEmail(email, otpCode); err != nil {
		log.Printf("❌ Ошибка при отправке OTP: %v", err)
		return fmt.Errorf("ошибка при отправке OTP: %v", err)
	}

	log.Printf("📨 OTP для сброса пароля отправлен: %s (OTP: %s)", email, otpCode)
	us.nats.Publish("password.reset.requested", []byte(fmt.Sprintf("Reset requested for Email: %s", email)))

	return nil
}

func (us *UserService) ResetPassword(email, otpCode, newPassword string) error {
	log.Printf("➡️ ResetPassword called: email=%s, otp=%s", email, otpCode)

	user, err := us.Repo.GetUserByEmail(email)
	if err != nil {
		log.Printf("❌ Пользователь не найден: %s", email)
		return fmt.Errorf("user not found")
	}

	storedOtp, ok := us.otps[user.ID]
	if !ok || storedOtp != otpCode {
		log.Printf("❌ Неверный OTP для пользователя: %s", user.ID)
		return fmt.Errorf("invalid OTP")
	}

	delete(us.otps, user.ID)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("❌ Ошибка при хэшировании нового пароля: %v", err)
		return fmt.Errorf("ошибка при хэшировании нового пароля: %v", err)
	}

	log.Printf("🔑 Пароль успешно обновлён для: %s", user.ID)
	us.nats.Publish("password.reset", []byte(fmt.Sprintf("Password reset for UserID: %s", user.ID)))

	return us.Repo.UpdatePassword(user.ID, string(hashedPassword))
}

func (us *UserService) GetUser(userID string) (*model.User, error) {
	log.Printf("➡️ Получение пользователя по ID: %s", userID)

	// Попытка получить пользователя из Redis
	cached, err := us.cache.Get(userID)
	if err == nil && cached != "" {
		var user model.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			log.Printf("⚡ Пользователь получен из Redis-кэша: %s", userID)
			return &user, nil
		} else {
			log.Printf("⚠️ Ошибка при десериализации пользователя из Redis: %v", err)
		}
	} else {
		log.Printf("ℹ️ Пользователь не найден в Redis или ошибка Redis: %v", err)
	}

	// Если нет в кэше — получаем из БД
	user, err := us.Repo.GetUserByID(userID)
	if err != nil {
		log.Printf("❌ Пользователь не найден в БД: %s", userID)
		return nil, err
	}

	// Сериализация пользователя в JSON
	userJSON, err := json.Marshal(user)
	if err != nil {
		log.Printf("⚠️ Ошибка при сериализации пользователя: %v", err)
		return nil, err
	}

	// Сохранение в Redis с TTL 5 минут
	if err := us.cache.Set(userID, string(userJSON), 5*time.Minute); err != nil {
		log.Printf("⚠️ Ошибка при сохранении пользователя в Redis: %v", err)
	}

	return user, nil
}



func (us *UserService) UpdateUser(userID, username, email string) error {
	log.Printf("✏️ Обновление пользователя: ID=%s, username=%s, email=%s", userID, username, email)

	// Удалим старую запись из кэша
	_ = us.cache.Del(userID)

	return us.Repo.UpdateUser(userID, username, email)
}

func (us *UserService) DeleteUser(userID string) error {
	log.Printf("🗑️ Удаление пользователя: %s", userID)

	// Удалим запись из Redis
	_ = us.cache.Del(userID)

	return us.Repo.DeleteUser(userID)
}

func (us *UserService) ListUsers() ([]model.User, error) {
	log.Printf("📋 Получение списка всех пользователей")
	return us.Repo.GetAllUsers()
}
