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

func (assigment AssignmentPsql) ReadByUserID(context.Context, entity.UserID) ([]*entity.Assignment, error) {
	return nil, nil
}

func (assigment AssignmentPsql) Save(context.Context, *entity.Assignment) error {
	return nil
}
func (assigment AssignmentPsql) Delete(context.Context, *entity.Assignment) error {
	return nil
}

func (assigment AssignmentPsql) Prune(ctx context.Context) error {
	query, args, err := assigment.pg.Builder().
		Delete(assigmentsTable).
		Where(sql.Eq{"deleted": true}).ToSql()
	if err != nil {
		// panic(err) ?
		return fmt.Errorf("failed buled query: %w", err)
	}
	tag, err := assigment.pg.Conn(ctx).Exec(ctx, query, args...)
	log.Info().Err(err).Msgf("prune %v assigments by %v", tag.RowsAffected(), tag.String())
	return err
}
