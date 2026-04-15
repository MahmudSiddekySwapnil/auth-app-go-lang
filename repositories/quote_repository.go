package repositories

import (
	"auth-app/models"
	"database/sql"
)

type QuoteRepository interface {
	Save(quote *models.Quote) error
}

type quoteRepository struct {
	db *sql.DB
}

func NewQuoteRepository(db *sql.DB) QuoteRepository {
	return &quoteRepository{db: db}
}

func (r *quoteRepository) Save(quote *models.Quote) error {
	err := r.db.QueryRow(
		"INSERT INTO quotes (name, quoteline) VALUES ($1, $2) RETURNING id",
		quote.Name,
		quote.QuoteLine,
	).Scan(&quote.ID)

	return err
}
