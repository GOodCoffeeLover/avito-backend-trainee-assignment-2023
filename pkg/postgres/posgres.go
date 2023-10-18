package postgres

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
)

func New(ctx context.Context, connString string) (*Postgres, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to psql: %w", err)
	}
	return &Postgres{
		pgConn: conn,
	}, nil

}

type Postgres struct {
	pgConn *pgx.Conn
}

func (p Postgres) Conn() *pgx.Conn {
	return p.pgConn
}

func (p Postgres) Close(ctx context.Context) error {
	return p.pgConn.Close(ctx)
}
