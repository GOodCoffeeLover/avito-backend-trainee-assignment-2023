package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/pkg/postgres"
	sql "github.com/Masterminds/squirrel"
	"github.com/rs/zerolog/log"
)

type AssignmentPsql struct {
	pg *postgres.Postgres
}

func NewAssignmentPsql(ctx context.Context, pg *postgres.Postgres) (*AssignmentPsql, error) {

	if pg == nil {
		return nil, fmt.Errorf("nil potgresql client")
	}
	assignments := &AssignmentPsql{
		pg: pg,
	}
	go func() {
		ctx := context.Background()
		for {
			<-time.After(120 * time.Second)
			log.Info().Err(assignments.Prune(ctx)).Msg("Pruned segments")
		}
	}()
	return assignments, nil
}

func (as AssignmentPsql) ReadByUserID(ctx context.Context, uid entity.UserID) ([]*entity.Assignment, error) {
	query, args, err := as.pg.Builder().
		Select(assignments.userID, assignments.segmentName).
		From(assignments.table).
		Where(sql.Eq{assignments.table + "." + assignments.userID: uid}).
		Where(sql.Eq{assignments.table + "." + assignments.deleted: false}).
		InnerJoin("( SELECT * FROM users WHERE deleted = FALSE ) AS users ON assignments.user_id = users.id").
		InnerJoin("( SELECT * FROM segments WHERE deleted = FALSE ) AS segments ON assignments.segment_name = segments.name").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed buled query: %w", err)
	}

	rows, err := as.pg.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed get users assignments: %w", err)
	}
	assignments := []*entity.Assignment{}
	for rows.Next() {
		var asssigment entity.Assignment
		err = rows.Scan(&asssigment.User, &asssigment.Segment)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignments: %w", err)
		}
		assignments = append(assignments, &asssigment)
	}
	return assignments, nil
}

func (as AssignmentPsql) Create(ctx context.Context, assignment *entity.Assignment) error {
	query, args, err := as.pg.Builder().
		Insert(assignments.table).
		Columns(assignments.userID, assignments.segmentName).
		Values(assignment.User, assignment.Segment).
		Suffix(
			fmt.Sprintf(
				`ON CONFLICT (%v, %v) DO UPDATE SET %v = ?, %v = ?, %v = FALSE WHERE %v = TRUE`,
				assignments.userID, assignments.segmentName,
				assignments.userID, assignments.segmentName, assignments.deleted,
				fmt.Sprintf("%v.%v", assignments.table, assignments.deleted),
			),
			assignment.User, assignment.Segment,
		).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed buled query: %w", err)
	}
	tag, err := as.pg.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed save assignment %v: %w", assignment, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("assigment %w", entity.ErrAlreadyExists)
	}
	return nil
}
func (as AssignmentPsql) Delete(ctx context.Context, assignment *entity.Assignment) error {
	query, args, err := as.pg.Builder().
		Update(assignments.table).
		Set(assignments.deleted, true).
		Where(sql.Eq{
			assignments.userID:      assignment.User,
			assignments.segmentName: assignment.Segment,
			assignments.deleted:     false,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("failed buled query: %w", err)
	}
	tag, err := as.pg.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed save assignment %v: %w", assignment, err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("failed save assignment %v: affected not 1 row but %v", assignment, tag.RowsAffected())
	}
	return nil
}

func (as AssignmentPsql) Prune(ctx context.Context) error {
	query, args, err := as.pg.Builder().
		Delete(assignments.table).
		Where(sql.Eq{assignments.deleted: true}).ToSql()
	if err != nil {
		return fmt.Errorf("failed buled query: %w", err)
	}
	tag, err := as.pg.Conn(ctx).Exec(ctx, query, args...)
	log.Info().Err(err).Msgf("prune %v assignments by %v", tag.RowsAffected(), tag.String())
	return err
}
