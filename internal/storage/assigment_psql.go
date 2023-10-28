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
	assigments := &AssignmentPsql{
		pg: pg,
	}
	go func() {
		ctx := context.Background()
		for {
			<-time.After(60 * time.Second)
			log.Info().Err(assigments.Prune(ctx)).Msg("Pruning segments")
		}
	}()
	return assigments, nil
}

func (as AssignmentPsql) ReadByUserID(ctx context.Context, uid entity.UserID) ([]*entity.Assignment, error) {
	// TODO: fix quering deleted users or segments
	query, args, err := as.pg.Builder().
		Select(assigmentsTable.userID, assigmentsTable.segmentName).
		From(assigmentsTable.name).
		Where(sql.Eq{assigmentsTable.userID: uid}).
		Where(sql.Eq{assigmentsTable.deleted: false}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed buled query: %w", err)
	}
	rows, err := as.pg.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed get users assigments: %w", err)
	}
	assigments := []*entity.Assignment{}
	for rows.Next() {
		var asssigment entity.Assignment
		err = rows.Scan(asssigment.User, asssigment.Segment)
		if err != nil {
			return nil, fmt.Errorf("failed to scan assigments %w", err)
		}
		assigments = append(assigments, &asssigment)
	}
	return assigments, nil
}

func (as AssignmentPsql) Create(ctx context.Context, assignment *entity.Assignment) error {
	updateQuery, updateArgs, err := as.pg.Builder().
		Update(assigmentsTable.name).
		Set(assigmentsTable.userID, assignment.User).
		Set(assigmentsTable.segmentName, assignment.Segment).
		Set(assigmentsTable.deleted, false).
		Where(sql.Eq{assigmentsTable.deleted: true}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed buled query: %w", err)
	}
	query, args, err := as.pg.Builder().
		Insert(assigmentsTable.name).
		Columns(assigmentsTable.userID, assigmentsTable.segmentName).
		Values(assignment.User, assignment.Segment).
		Suffix(fmt.Sprintf("ON CONFLICT (%v, %v) DO %v",
			assigmentsTable.userID, assigmentsTable.segmentName, updateQuery), updateArgs...).
		ToSql()
	fmt.Println(query)
	if err != nil {
		return fmt.Errorf("failed buled query: %w", err)
	}
	_, err = as.pg.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed save assigment %v: %w", assignment, err)
	}
	return nil
}
func (as AssignmentPsql) Delete(ctx context.Context, assignment *entity.Assignment) error {
	query, args, err := as.pg.Builder().
		Update(assigmentsTable.name).
		Set(assigmentsTable.deleted, true).
		Where(sql.Eq{
			assigmentsTable.userID:      assignment.User,
			assigmentsTable.segmentName: assignment.Segment,
			assigmentsTable.deleted:     false,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("failed buled query: %w", err)
	}
	tag, err := as.pg.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed save assigment %v: %w", assignment, err)
	}
	if tag.RowsAffected() != 1 {
		return fmt.Errorf("failed save assigment %v: affected not 1 row but %v", assignment, tag.RowsAffected())
	}
	return nil
}

func (as AssignmentPsql) Prune(ctx context.Context) error {
	query, args, err := as.pg.Builder().
		Delete(assigmentsTable.name).
		Where(sql.Eq{assigmentsTable.deleted: true}).ToSql()
	if err != nil {
		return fmt.Errorf("failed buled query: %w", err)
	}
	tag, err := as.pg.Conn(ctx).Exec(ctx, query, args...)
	log.Info().Err(err).Msgf("prune %v assigments by %v", tag.RowsAffected(), tag.String())
	return err
}
