package models

import (
	"time"
)

type URLShortener struct {
	ID          int       `db:"id"`
	OriginalURL string    `db:"original_url"`
	ShortCode   string    `db:"short_code"`
	CreatedAt   time.Time `db:"created_at"`
	ExpiredAt   time.Time `db:"expired_at"`
	Status      bool      `db:"status"`
}
