package repository

import (
	"context"
	"fmt"

	workspacedb "github.com/findardi/Wadi/server/internal/workspace/repository/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	*workspacedb.Queries
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Queries: workspacedb.New(pool),
		pool:    pool,
	}
}

func (r *Repository) ExecTx(ctx context.Context, fn func(*workspacedb.Queries) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := fn(r.Queries.WithTx(tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
