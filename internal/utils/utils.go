package utils

import (
	"crypto/md5"
	"sort"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	// Хэшируем пароль с коэффициентом сложности 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// GenerateDialogKey принимает два UUID пользователя и возвращает UUID диалога
func GenerateDialogKey(user1, user2 string) string {
	// сортируем uuid как строки, чтобы порядок не влиял
	users := []string{user1, user2}
	sort.Strings(users)

	// склеиваем и считаем md5
	h := md5.Sum([]byte(users[0] + users[1]))

	// md5 это 16 байт — как раз подходит под uuid
	dialogKey, err := uuid.FromBytes(h[:])
	if err != nil {
		// не должен падать, но на всякий случай
		panic(err)
	}
	return dialogKey.String()
}
