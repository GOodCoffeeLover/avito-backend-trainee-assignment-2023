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
			<-time.After(60 * time.Second)
			log.Info().Err(assignments.Prune(ctx)).Msg("Pruning segments")
		}
	}()
	return assignments, nil
}

func (as AssignmentPsql) ReadByUserID(ctx context.Context, uid entity.UserID) ([]*entity.Assignment, error) {
	// TODO: fix quering deleted users or segments
	query, args, err := as.pg.Builder().
		Select(assignmentsTable.userID, assignmentsTable.segmentName).
		From(assignmentsTable.name).
		Where(sql.Eq{assignmentsTable.userID: uid}).
		Where(sql.Eq{assignmentsTable.deleted: false}).
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
		err = rows.Scan(asssigment.User, asssigment.Segment)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assignments %w", err)
		}
		assignments = append(assignments, &asssigment)
	}
	return assignments, nil
}

func (as AssignmentPsql) Create(ctx context.Context, assignment *entity.Assignment) error {
	updateQuery, updateArgs, err := as.pg.Builder().
		Update(assignmentsTable.name).
		Set(assignmentsTable.userID, assignment.User).
		Set(assignmentsTable.segmentName, assignment.Segment).
		Set(assignmentsTable.deleted, false).
		Where(sql.Eq{assignmentsTable.deleted: true}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed buled query: %w", err)
	}
	query, args, err := as.pg.Builder().
		Insert(assignmentsTable.name).
		Columns(assignmentsTable.userID, assignmentsTable.segmentName).
		Values(assignment.User, assignment.Segment).
		Suffix(fmt.Sprintf("ON CONFLICT (%v, %v) DO %v",
			assignmentsTable.userID, assignmentsTable.segmentName, updateQuery), updateArgs...).
		ToSql()
	fmt.Println(query)
	if err != nil {
		return fmt.Errorf("failed buled query: %w", err)
	}
	_, err = as.pg.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed save assignment %v: %w", assignment, err)
	}
	return nil
}
func (as AssignmentPsql) Delete(ctx context.Context, assignment *entity.Assignment) error {
	query, args, err := as.pg.Builder().
		Update(assignmentsTable.name).
		Set(assignmentsTable.deleted, true).
		Where(sql.Eq{
			assignmentsTable.userID:      assignment.User,
			assignmentsTable.segmentName: assignment.Segment,
			assignmentsTable.deleted:     false,
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
		Delete(assignmentsTable.name).
		Where(sql.Eq{assignmentsTable.deleted: true}).ToSql()
	if err != nil {
		return fmt.Errorf("failed buled query: %w", err)
	}
	tag, err := as.pg.Conn(ctx).Exec(ctx, query, args...)
	log.Info().Err(err).Msgf("prune %v assignments by %v", tag.RowsAffected(), tag.String())
	return err
}
