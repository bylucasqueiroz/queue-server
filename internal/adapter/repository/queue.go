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

type PostgresQueueRepository struct {
	db *sql.DB
}

func NewPostgresQueueRepository(config *config.Config) (*PostgresQueueRepository, error) {
	db, err := sql.Open("postgres", config.ConString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &PostgresQueueRepository{db: db}, nil
}

func (r *PostgresQueueRepository) Save(ctx context.Context, queue *domain.Queue) error {
	query := `INSERT INTO queues (name, created_at) VALUES ($1, $2) RETURNING name`
	err := r.db.QueryRowContext(ctx, query, queue.Name, time.Now()).Scan(&queue.Name)
	if err != nil {
		return fmt.Errorf("failed to save queue: %v", err)
	}
	return nil
}

func (r *PostgresQueueRepository) GetByName(ctx context.Context, name string) (*domain.Queue, error) {
	query := `SELECT name, created_at FROM queues WHERE name = $1`
	row := r.db.QueryRowContext(ctx, query, name)

	queue := &domain.Queue{}
	if err := row.Scan(&queue.Name, &queue.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get queue: %v", err)
	}
	return queue, nil
}

func (r *PostgresQueueRepository) Delete(ctx context.Context, name string) error {
	query := `DELETE FROM queues WHERE name = $1`
	_, err := r.db.ExecContext(ctx, query, name)
	if err != nil {
		return fmt.Errorf("failed to delete queue: %v", err)
	}
	return nil
}
