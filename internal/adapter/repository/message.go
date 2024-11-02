package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"queueserver/internal/adapter/config"
	"queueserver/internal/core/domain"

	_ "github.com/lib/pq"
)

type PostgresMessageRepository struct {
	db *sql.DB
}

func NewPostgresMessageRepository(config *config.Config) (*PostgresMessageRepository, error) {
	db, err := sql.Open("postgres", config.ConString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &PostgresMessageRepository{db: db}, nil
}

func (r *PostgresMessageRepository) Save(ctx context.Context, message *domain.Message) error {
	query := `INSERT INTO messages (id, body, receipt_handle, visibility_timeout, queue_name) 
              VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO UPDATE 
              SET body = EXCLUDED.body, receipt_handle = EXCLUDED.receipt_handle, 
                  visibility_timeout = EXCLUDED.visibility_timeout, queue_name = EXCLUDED.queue_name`
	_, err := r.db.ExecContext(ctx, query, message.ID, message.Body, message.ReceiptHandle, message.VisibilityTimeout, message.QueueName)
	if err != nil {
		return fmt.Errorf("failed to save message: %v", err)
	}
	return nil
}

func (r *PostgresMessageRepository) GetByID(ctx context.Context, id string) (*domain.Message, error) {
	query := `SELECT id, body, receipt_handle, visibility_timeout, queue_name FROM messages WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	message := &domain.Message{}
	if err := row.Scan(&message.ID, &message.Body, &message.ReceiptHandle, &message.VisibilityTimeout, &message.QueueName); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get message: %v", err)
	}
	return message, nil
}

func (r *PostgresMessageRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM messages WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %v", err)
	}
	return nil
}
