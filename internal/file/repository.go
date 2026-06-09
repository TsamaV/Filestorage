package file

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct {
	db *pgxpool.Pool
}

func NewFileRepository(database *pgxpool.Pool) *FileRepository {
	return &FileRepository{
		db: database,
	}
}

func (repo *FileRepository) Create(ctx context.Context, userID int, url, fileName string, size int) (int, error) {
	var id int

	query := `
	INSERT INTO files (user_id, url, file_name, size, created_at)
	VALUES($1, $2, $3, $4, $5)
	RETURNING id
	`

	err := repo.db.QueryRow(ctx, query, userID, url, fileName, size, time.Now()).Scan(&id)

	return id, err
}

func (repo *FileRepository) GetFiles(ctx context.Context, page, limit int) ([]File, error) {
	var result []File

	query := `
	SELECT * FROM files
	ORDER BY id 
	LIMIT $1
	OFFSET $2
	`

	offset := (page - 1) * limit

	rows, err := repo.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var f File

		err := rows.Scan(&f.ID, &f.UserID, &f.URL, &f.FileName, &f.Size, &f.CreatedAt)
		if err != nil {
			return nil, err
		}

		result = append(result, f)
	}

	return result, err
}

func (repo *FileRepository) GetByID(ctx context.Context, id int) (File, error) {
	var f File

	query := `
	SELECT * FROM files
	WHERE id = $1
	`

	err := repo.db.QueryRow(ctx, query, id).Scan(&f.ID, &f.UserID, &f.URL, &f.FileName, &f.Size, &f.CreatedAt)

	return f, err
}

func (repo *FileRepository) Delete(ctx context.Context, id int) error {
	query := `
	DELETE FROM files
	WHERE id = $1
	`

	_, err := repo.db.Exec(ctx, query, id)

	return err
}
