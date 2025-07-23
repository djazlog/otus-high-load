package main

import (
	"bufio"
	"database/sql"
	"runtime"

	"fmt"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"

	"log"
	"math/rand"
	"os"
	"otus-project/internal/config"
	"strings"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const maxWorkers = 10 // количество горутин для вставки
const chunkSize = 500 // количество строк на один чанк для чтения

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

func main() {
	fmt.Println(runtime.NumCPU())
	password, _ := hashPassword("123")
	runtime.GOMAXPROCS(6 * runtime.NumCPU())
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

	userChan := make(chan []User, chunkSize*maxWorkers)
	var wg sync.WaitGroup

	// Запускаем воркеры для вставки
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for users := range userChan {
				log.Printf("Воркер %d: Получено %d пользователей", i, len(users))
				if err := batchInsert(db, users); err != nil {
					log.Printf("Ошибка при батчевой вставке: %v", err)
				} else {
					log.Printf("Добавлено %d пользователей", len(users))
				}
			}
		}()
	}

	// Параллельное чтение файла
	readFileInChunks(file, userChan, password)
	time.Sleep(5 * time.Second) // дать время на обработку
	close(userChan)
	wg.Wait()
	log.Println("Загрузка завершена.")
}

// readFileInChunks - разбивает файл на чанки и обрабатывает каждый в отдельной горутине
func readFileInChunks(file *os.File, userChan chan<- []User, password string) {

	log.Println("Начинаем чтение файла")
	scanner := bufio.NewScanner(file)
	lines := []string{}

	// Считываем все строки в память
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Удаляем заголовок
	if len(lines) > 0 && strings.HasPrefix(lines[0], "first_name last_name") {
		lines = lines[1:]
	}

	log.Printf("Чтение файла завершено. Всего строк: %d", len(lines))

	// Разбиваем на чанки
	var wgRead sync.WaitGroup
	for i := 0; i < len(lines); i += chunkSize {
		end := i + chunkSize
		if end > len(lines) {
			end = len(lines)
		}
		chunk := lines[i:end]
		wgRead.Add(1)
		go func(chunk []string) {
			defer wgRead.Done()

			batch := parseChunk(chunk, password)
			log.Printf("Готово %d строк, отправляем в поток", len(batch))
			userChan <- batch
		}(chunk)
	}
	wgRead.Wait()
}

// parseChunk - парсит чанк строк и формирует батч пользователей
func parseChunk(lines []string, password string) []User {
	var batch []User

	for _, line := range lines {
		record := strings.Split(line, ",")
		if len(record) < 3 {
			continue
		}

		fullName := record[0]
		nameParts := strings.Fields(fullName)
		if len(nameParts) < 2 {
			continue
		}

		firstName := nameParts[0]
		lastName := nameParts[1]

		birthDate, err := time.Parse("2006-01-02", record[1])
		if err != nil {
			continue
		}

		city := record[2]
		id, _ := uuid.NewV4()

		user := User{
			Id:        id,
			FirstName: firstName,
			LastName:  lastName,
			BirthDate: birthDate,
			City:      city,
			Password:  password,
			Biography: randomString(128),
			CreateAt:  time.Now(),
		}

		batch = append(batch, user)
	}
	log.Printf("Парсинг чанка завершен. Всего пользователей: %d", len(batch))
	return batch
}

// batchInsert - вставляет батч пользователей в БД
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

func hashPassword(password string) (string, error) {
	// Хэшируем пароль с коэффициентом сложности 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
