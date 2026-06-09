package file

import "time"

type File struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	URL       string    `json:"url"`
	FileName  string    `json:"file_name"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}