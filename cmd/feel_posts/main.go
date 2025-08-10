package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"os"
	"otus-project/internal/config"
	"runtime"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {

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
	file, err := os.Open("./docs/posts.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Читаем все строки из файла
	scanner := bufio.NewScanner(file)
	var posts []string
	for scanner.Scan() {
		posts = append(posts, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Получаем список пользователей
	rows, err := db.Query("SELECT id FROM users ORDER BY random() Limit 10000 ")
	if err != nil {
		log.Fatal(err)
	}

	var userIDs []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			log.Fatal(err)
		}
		userIDs = append(userIDs, userID)
	}

	// Вставляем 10 постов для каждого пользователя
	for _, userID := range userIDs {
		for i := 0; i < 10; i++ {
			if len(posts) == 0 {
				break
			}
			//TODO: Убрать, сейчас для теста на одного опльзователя сохраняем
			userID = "d77da351-b954-4f72-b51b-d94e97c13fc9"

			content := posts[0]
			posts = posts[1:]
			id := uuid.New().String()

			_, err := db.Exec("INSERT INTO posts (id, author_user_id, content, created_at) VALUES ($1, $2, $3, $4)", id, userID, content, time.Now())
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return
	//Генерация случайных связей между пользователями
	q := `
		WITH user_ids AS (
				SELECT id
				FROM users
				WHERE id IN (%s)
			),
			random_friends AS (
				SELECT DISTINCT ON (u1.id)
					u1.id AS user_id,
					u2.id AS friend_id
				FROM user_ids u1
				CROSS JOIN LATERAL (
					SELECT id
					FROM users u2
					WHERE u2.id != u1.id
					  AND NOT EXISTS (
						  SELECT 1
						  FROM friends f
						  WHERE f.user_id = u1.id
							AND f.friend_id = u2.id
					  )
					ORDER BY random()
					LIMIT 10
				) u2
				ORDER BY u1.id, random()
			)
			INSERT INTO friends (user_id, friend_id)
			SELECT user_id, friend_id
			FROM random_friends;
`
	userIDPlaceholders := make([]string, len(userIDs))
	for i := range userIDs {
		userIDPlaceholders[i] = fmt.Sprintf("'%s'", userIDs[i])
	}
	queryWithIDs := fmt.Sprintf(q, strings.Join(userIDPlaceholders, ", "))

	_, err = db.Exec(queryWithIDs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Посты успешно импортированы!")
}
