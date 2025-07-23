package user

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"otus-project/internal/client/db"
	"otus-project/internal/model"
	"otus-project/internal/repository"
	"otus-project/internal/repository/user/converter"
	modelRepo "otus-project/internal/repository/user/model"
	"time"
)

const (
	tableName = "users"

	idColumn         = "id"
	firstNameColumn  = "first_name"
	secondNameColumn = "second_name"
	birthDateColumn  = "birth_date"
	biographyColumn  = "biography"
	cityColumn       = "city"
	passwordColumn   = "password"
	createdAtColumn  = "created_at"
	updatedAtColumn  = "updated_at"
)

type repo struct {
	db db.Client
}

// Секретный ключ, который используется для подписи токена
// TODO: вынести в конфиг
var jwtSecret = []byte("my-super-secret-key")

func generateJWT(userID string) (string, error) {
	// Создаём новый токен с алгоритмом HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // истекает через сутки
		"iat":     time.Now().Unix(),
	})

	// Подписываем токен
	return token.SignedString(jwtSecret)
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

// Login вход пользователя.
func (r *repo) Login(ctx context.Context, login *model.LoginDto) (*string, error) {
	builder := sq.Select(idColumn, passwordColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: login.Id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Login",
		QueryRaw: query,
	}

	var userId, hashedPassword string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userId, &hashedPassword)
	if err != nil {
		return nil, err
	}

	ok := checkPasswordHash(login.Password, hashedPassword)

	if !ok {
		return nil, errors.New("invalid password")
	}

	jwtNew, err := generateJWT(userId)
	if err != nil {
		return nil, err
	}

	return &jwtNew, nil
}

func checkPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// Register регистрация пользователя.
func (r *repo) Register(ctx context.Context, info *model.UserInfo) (string, error) {
	idNew, err := uuid.NewV4()
	if err != nil {
		return "", err
	}
	if info.Password == nil {
		return "", errors.New("password is required")
	}

	hashedPassword, err := hashPassword(*info.Password)
	if err != nil {
		return "", err
	}

	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(idColumn, firstNameColumn, secondNameColumn, birthDateColumn, biographyColumn, cityColumn, createdAtColumn, updatedAtColumn, passwordColumn).
		Values(idNew.String(), info.FirstName, info.SecondName, info.Birthdate, info.Biography, info.City, time.Now(), time.Now(), hashedPassword).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()
	log.Println(query, args, err)
	if err != nil {
		return "", err
	}

	q := db.Query{
		Name:     "user_repository.Register",
		QueryRaw: query,
	}

	var id string
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func hashPassword(password string) (string, error) {
	// Хэшируем пароль с коэффициентом сложности 10
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Get получение информации о пользователе по id.
func (r *repo) Get(ctx context.Context, id string) (*model.UserInfo, error) {
	builder := sq.Select(idColumn, firstNameColumn, secondNameColumn, birthDateColumn, biographyColumn, cityColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Get",
		QueryRaw: query,
	}

	var user modelRepo.User
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return converter.ToUserInfoFromRepo(&user), nil
}

// Search поиск пользователей.
func (r *repo) Search(ctx context.Context, filter *model.UserFilter) ([]*model.UserInfo, error) {
	builder := sq.Select(idColumn, firstNameColumn, secondNameColumn, birthDateColumn, biographyColumn, cityColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName)

	// Фильтрация по firstName и secondName
	if filter.FirstName != "" {
		builder = builder.Where(sq.Like{firstNameColumn: filter.FirstName + "%"})
	}

	if filter.LastName != "" {
		builder = builder.Where(sq.Like{secondNameColumn: filter.LastName + "%"})
	}

	builder = builder.OrderBy(idColumn).Limit(300)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "user_repository.Search",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.UserInfo
	for rows.Next() {
		var user modelRepo.User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.SecondName, &user.Birthdate, &user.Biography, &user.City, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, converter.ToUserInfoFromRepo(&user))
	}

	return users, nil
}
