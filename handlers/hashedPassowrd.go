package handlers

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func HashPassword(password string) string {
	// Генерируем хэш пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Ошибка при хэшировании пароля:", err)
		return "" // Возвращаем ошибку, если произошла ошибка
	}
	return string(hashedPassword) // Возвращаем хэшированный пароль
}

// CheckPasswordHash проверяет совпадение пароля с хэшем
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err != nil // Если ошибка равна nil, значит пароль совпадает
}
