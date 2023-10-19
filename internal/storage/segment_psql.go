package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type SegmentPsqlStorage struct {
	pg *postgres.Postgres
}

func NewSegmentPsqlStorage(ctx context.Context, pg *postgres.Postgres) (*SegmentPsqlStorage, error) {
	if pg == nil {
		return nil, fmt.Errorf("nil potgresql client")
	}
	segStorage := &SegmentPsqlStorage{
		pg: pg,
	}
	go func() {
		ctx := context.Background()
		for {
			<-time.After(60 * time.Second)
			log.Info().Err(segStorage.Prune(ctx)).Msg("Pruning segments")
		}
	}()
	return segStorage, nil
}

func (sps SegmentPsqlStorage) Create(ctx context.Context, segment *entity.Segment) error {
	// FIXME: conflict when deleted but not pruned
	query := `
	INSERT INTO segments (name) VALUES ($1)
	ON CONFLICT (name) 
	DO UPDATE 
		SET deleted=FALSE
		WHERE segments.deleted=TRUE`
	tag, err := sps.pg.Conn(ctx).Exec(ctx, query, segment.Name)
	if err != nil {
		return fmt.Errorf("failed create segment %v: %w", segment.Name, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed create segment %v: %w", segment.Name, ErrAlreadyExists)
	}
	return nil
}

func (sps SegmentPsqlStorage) ReadByName(ctx context.Context, segmentName entity.SegmentName) (*entity.Segment, error) {
	row := sps.pg.Conn(ctx).QueryRow(ctx, "SELECT name FROM segments WHERE name=$1 AND deleted=FALSE", segmentName)

	segment := &entity.Segment{}
	if err := row.Scan(&segment.Name); err != nil {
		return nil, fmt.Errorf("failed read segment by name %v: %w", segmentName, err)
	}
	return segment, nil
}
func (sps SegmentPsqlStorage) ReadAll(ctx context.Context) ([]*entity.Segment, error) {
	rows, err := sps.pg.Conn(ctx).Query(ctx, "SELECT name FROM segments WHERE deleted=FALSE")
	if err != nil {
		return nil, fmt.Errorf("failed read all segments: %w", err)
	}

	segments := []*entity.Segment{}
	for rows.Next() {
		segment := &entity.Segment{}
		if err := rows.Scan(&segment.Name); err != nil {
			return nil, fmt.Errorf("failed scan segment: %w", err)
		}
		segments = append(segments, segment)
	}
	return segments, nil
}

func (sps SegmentPsqlStorage) Delete(ctx context.Context, segmentName entity.SegmentName) error {
	tag, err := sps.pg.Conn(ctx).Exec(ctx, "UPDATE segments SET deleted=TRUE WHERE name=$1 AND deleted=FALSE", segmentName)
	if err != nil {
		return fmt.Errorf("failed deleting segment %v: %w", segmentName, err)
	}
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (sps SegmentPsqlStorage) Prune(ctx context.Context) error {
	if _, err := sps.pg.Conn(ctx).Exec(ctx, "DELETE FROM segments WHERE deleted=TRUE"); err != nil {
		return fmt.Errorf("failed prune segments: %w", err)
	}
	return nil
}
