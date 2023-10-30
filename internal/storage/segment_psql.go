package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type SegmentPsql struct {
	pg *postgres.Postgres
}

func NewSegmentPsql(ctx context.Context, pg *postgres.Postgres) (*SegmentPsql, error) {
	if pg == nil {
		return nil, fmt.Errorf("nil potgresql client")
	}
	segments := &SegmentPsql{
		pg: pg,
	}
	go func() {
		ctx := context.Background()
		for {
			<-time.After(60 * time.Second)
			log.Info().Err(segments.Prune(ctx)).Msg("Pruning segments")
		}
	}()
	return segments, nil
}

func (sps SegmentPsql) Create(ctx context.Context, segment *entity.Segment) error {
	query := `
    INSERT INTO segments (name) VALUES ($1)
    ON CONFLICT (name) 
    DO UPDATE 
        SET deleted=FALSE
        WHERE segments.deleted=TRUE`
	tag, err := sps.pg.Conn(ctx).Exec(ctx, query, segment.Name)
	if err != nil {
		return fmt.Errorf("failed create segment %+v: %w", segment, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed create segment %+v: %w", segment, ErrAlreadyExists)
	}
	return nil
}

func (sps SegmentPsql) ReadByName(ctx context.Context, segmentName entity.SegmentName) (*entity.Segment, error) {
	row := sps.pg.Conn(ctx).QueryRow(ctx, "SELECT name FROM segments WHERE name=$1 AND deleted=FALSE", segmentName)

	segment := &entity.Segment{}
	if err := row.Scan(&segment.Name); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, fmt.Errorf("segmnet %v not found: %w", segmentName, ErrNotFound)
		default:
			return nil, fmt.Errorf("failed scan segment by name %v: %w", segmentName, err)
		}
	}
	return segment, nil
}
func (sps SegmentPsql) ReadAll(ctx context.Context) ([]*entity.Segment, error) {
	rows, err := sps.pg.Conn(ctx).Query(ctx, "SELECT name FROM segments WHERE deleted=FALSE")
	if err != nil {
		return nil, fmt.Errorf("failed read all segments: %w", err)
	}

	segments := []*entity.Segment{}
	for rows.Next() {
		segment := &entity.Segment{}
		if err := rows.Scan(&segment.Name); err != nil {
			return nil, fmt.Errorf("failed scan segment row: %w", err)
		}
		segments = append(segments, segment)
	}
	return segments, nil
}

func (sps SegmentPsql) Delete(ctx context.Context, segmentName entity.SegmentName) error {
	tag, err := sps.pg.Conn(ctx).Exec(ctx, "UPDATE segments SET deleted=TRUE WHERE name=$1 AND deleted=FALSE", segmentName)
	if err != nil {
		return fmt.Errorf("failed deleting segment %v: %w", segmentName, err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (sps SegmentPsql) Prune(ctx context.Context) error {
	tag, err := sps.pg.Conn(ctx).Exec(ctx, "DELETE FROM segments WHERE deleted=TRUE")
	if err != nil {
		return fmt.Errorf("failed prune segments: %w", err)
	}
	log.Info().Msgf("prune %v segments by %v", tag.RowsAffected(), tag.String())
	return nil
}
