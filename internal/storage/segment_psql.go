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
	if _, err := sps.pgConn.Exec(ctx, "INSERT INTO segments (name) VALUES ($1)", segment.Name); err != nil {
		return fmt.Errorf("failed create segment %v: %w", segment.Name, err)
	}
	return nil
}

func (sps SegmentPsqlStorage) GetByName(ctx context.Context, segmentName entity.SegmentName) (*entity.Segment, error) {
	row := sps.pgConn.QueryRow(ctx, "SELECT (name) FROM segments WHERE name=$1", segmentName)
	segment := &entity.Segment{}
	if err := row.Scan(&segment.Name); err != nil {
		return nil, err
	}
	return segment, nil
}
