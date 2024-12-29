package repositories

import (
	"github.com/fikrihkll/chat-app/application/chat"
	"context"
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserRepositoryPostgree struct {
	db *sql.DB
}

func NewUserRepositoryPostgree(db *sql.DB) chat.IUserRepository {
	return &UserRepositoryPostgree{db}
}

func (repo *UserRepositoryPostgree) GetUserByID(ctx context.Context, id string) (user chat.User, err error) {
	sql := "SELECT * FROM users WHERE id = $1"
	row := repo.db.QueryRowContext(ctx, sql, id)
	if row.Err() != nil {
		err = row.Err()
		log.Fatalf("Failed to select users: %v", err)
		return
	}

	if err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return
	}

	return
}

func (repo *UserRepositoryPostgree) GetUserByEmail(ctx context.Context, email string) (user chat.User, err error) {
	sql := "SELECT * FROM users WHERE email = $1"
	row := repo.db.QueryRowContext(ctx, sql, email)
	if row.Err() != nil {
		err = row.Err()
		log.Fatalf("Failed to select users: %v", err)
		return
	}

	if err = row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return
	}

	return
}

func (repo *UserRepositoryPostgree) CreateUser(ctx context.Context, user chat.User) (err error){
	sql := "INSERT INTO users VALUES(uuid_generate_v4(), $1, $2, $3, NOW(), NOW())"

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		return
	}

	_, sqlErr := repo.db.Exec(sql, user.Name, user.Email, hash)
	if sqlErr != nil {
		err = sqlErr
		log.Fatalf("Failed to insert users: %v", err)
		return
	}

	return
}