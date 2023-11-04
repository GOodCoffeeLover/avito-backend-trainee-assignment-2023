package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/pkg/postgres"
	sql "github.com/Masterminds/squirrel"
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
	query, args, err := sps.pg.Builder().
		Insert(segments.table).
		Columns(segments.name).
		Values(segment.Name).
		Suffix(fmt.Sprintf(`
		ON CONFLICT (%v) DO UPDATE 
		SET %v=?,
			%v=FALSE
		WHERE %v.%v=TRUE`,
			segments.name,
			segments.name,
			segments.deleted,
			segments.table, segments.deleted,
		),
			segment.Name,
		).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	tag, err := sps.pg.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed create segment %+v: %w", segment, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed create segment %+v: %w", segment, entity.ErrAlreadyExists)
	}
	return nil
}

func (sps SegmentPsql) ReadByName(ctx context.Context, segmentName entity.SegmentName) (*entity.Segment, error) {
	query, args, err := sps.pg.Builder().
		Select(segments.name).
		From(segments.table).
		Where(sql.Eq{
			segments.name:    segmentName,
			segments.deleted: false,
		}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	row := sps.pg.Conn(ctx).QueryRow(ctx, query, args...)

	segment := &entity.Segment{}
	if err := row.Scan(&segment.Name); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, fmt.Errorf("segmnet %v %w", segmentName, entity.ErrNotFound)
		default:
			return nil, fmt.Errorf("failed scan segment by name %v: %w", segmentName, err)
		}
	}
	return segment, nil
}
func (sps SegmentPsql) ReadAll(ctx context.Context) ([]*entity.Segment, error) {
	query, args, err := sps.pg.Builder().
		Select(segments.name).
		From(segments.table).
		Where(sql.Eq{
			segments.deleted: false,
		}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := sps.pg.Conn(ctx).Query(ctx, query, args...)
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
	query, args, err := sps.pg.Builder().
		Update(segments.table).
		Set(segments.deleted, true).
		Where(sql.Eq{
			segments.name:    segmentName,
			segments.deleted: false,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	tag, err := sps.pg.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed deleting segment %v: %w", segmentName, err)
	}
	if tag.RowsAffected() == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (sps SegmentPsql) Prune(ctx context.Context) error {
	query, args, err := sps.pg.Builder().
		Delete(segments.table).
		Where(sql.Eq{
			segments.deleted: true,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	tag, err := sps.pg.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed prune segments: %w", err)
	}
	log.Info().Msgf("prune %v segments by %v", tag.RowsAffected(), tag.String())
	return nil
}
