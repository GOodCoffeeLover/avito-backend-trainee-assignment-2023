package storage

import (
	"context"
	"fmt"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/jackc/pgx/v5"
)

type SegmentPsqlStorage struct {
	pgConn *pgx.Conn
}

func NewSegmentPsqlStorage(ctx context.Context, connStrign string) (*SegmentPsqlStorage, error) {

	conn, err := pgx.Connect(ctx, connStrign)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to psql: %w", err)
	}
	return &SegmentPsqlStorage{
		pgConn: conn,
	}, nil
}

func (sps SegmentPsqlStorage) CloseConn(ctx context.Context) error {
	return sps.pgConn.Close(ctx)

}
func (sps SegmentPsqlStorage) Create(ctx context.Context, segment entity.Segment) error {
	err := sps.pgConn.QueryRow(ctx, "INSERT INTO segments (name) VALUES ($1)", segment.Name()).Scan()
	return err
}
