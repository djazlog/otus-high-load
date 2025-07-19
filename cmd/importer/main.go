package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"os"
	"otus-project/internal/config"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib" // Драйвер PostgreSQL
)

type User struct {
	Id        uuid.UUID
	FirstName string
	LastName  string
	BirthDate time.Time
	City      string
	Password  string
	Biography string
	CreateAt  time.Time
}

const batchSize = 1000 // Количество записей на один батч

func main() {
	err := config.Load(".env")
	if err != nil {
		log.Fatal(fmt.Sprintf("Ошибка при получении env: %v", err))
	}

	dns, err := config.NewPGConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("Ошибка при получении конфигурации: %v", err))
	}

	// Подключение к БД
	db, err := sql.Open("pgx", dns.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Открываем CSV файл
	file, err := os.Open("./docs/people.v2.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Пропускаем заголовок
	var users []User
	for i, record := range records {
		if i == 0 {
			continue
		}

		fullName := record[0]
		nameParts := parseFullName(fullName)
		if len(nameParts) < 2 {
			log.Printf("Ошибка в строке: %v", record)
			continue
		}

		firstName := nameParts[0]
		lastName := nameParts[1]

		birthDate, err := time.Parse("2006-01-02", record[1])
		if err != nil {
			log.Printf("Не удалось распарсить дату рождения: %s", record[1])
			continue
		}

		city := record[2]
		id, _ := uuid.NewV4()
		password, _ := hashPassword("123")

		user := User{
			Id:        id,
			FirstName: firstName,
			LastName:  lastName,
			BirthDate: birthDate,
			City:      city,
			Password:  password,
			Biography: randomString(256),
			CreateAt:  time.Now(),
		}

		users = append(users, user)

		// Выполняем батчевую вставку
		if len(users) >= batchSize || i == len(records)-1 {
			if err := batchInsert(db, users); err != nil {
				log.Printf("Ошибка при батчевой вставке: %v", err)
			} else {
				log.Printf("Добавлено %d пользователей", len(users))
			}
			users = users[:0] // Очистка среза для следующего батча
		}
		log.Printf("Обработано %d строк", i)
	}
}

// Функция для разделения полного имени на имя и фамилию
func parseFullName(fullName string) []string {
	return strings.Fields(fullName)
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Генерация случайной строки длины n
func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Батчевая вставка пользователей
func batchInsert(db *sql.DB, users []User) error {
	if len(users) == 0 {
		return nil
	}

	query := `
		INSERT INTO users (
			id, first_name, second_name, birth_date, city, password, biography, created_at
		) VALUES `

	args := make([]interface{}, 0, len(users)*8)
	for i, user := range users {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", i*8+1, i*8+2, i*8+3, i*8+4, i*8+5, i*8+6, i*8+7, i*8+8)

		args = append(args,
			user.Id,
			user.FirstName,
			user.LastName,
			user.BirthDate,
			user.City,
			user.Password,
			user.Biography,
			user.CreateAt,
		)
	}

	_, err := db.Exec(query, args...)
	return err
}
