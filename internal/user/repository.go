package user

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(database *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: database,
	}
}

func (repo *UserRepository) CreateUser(ctx context.Context, email string, passwordHash string) (int, error) {
	var id int

	query := `
	INSERT INTO users (email, password_hash)
	VALUES($1, $2)
	RETURNING id
	`

	err := repo.db.QueryRow(ctx, query, email, passwordHash).Scan(&id)

	return id, err
}

func (repo *UserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var user User

	query := `
	SELECT * FROM users
	WHERE email = $1
	`

	err := repo.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, err
}
