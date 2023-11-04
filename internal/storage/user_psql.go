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

type UserPsql struct {
	pg *postgres.Postgres
}

func NewUserPsql(ctx context.Context, pg *postgres.Postgres) (*UserPsql, error) {
	if pg == nil {
		return nil, fmt.Errorf("nil potgresql client")
	}
	users := &UserPsql{
		pg: pg,
	}
	go func() {
		ctx := context.Background()
		for {
			<-time.After(120 * time.Second)
			log.Info().Err(users.Prune(ctx)).Msg("Pruned users")
		}
	}()
	return users, nil
}

func (sp UserPsql) Create(ctx context.Context, user *entity.User) error {
	query, args, err := sp.pg.Builder().
		Insert(users.table).
		Columns(users.id).
		Values(user.ID).
		Suffix(fmt.Sprintf(`
		ON CONFLICT (%v) DO UPDATE 
		SET %v=?,
			%v=FALSE
		WHERE %v.%v=TRUE`,
			users.id,
			users.id,
			users.deleted,
			users.table, segments.deleted,
		),
			user.ID,
		).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	tag, err := sp.pg.Conn(ctx).Exec(ctx, query, args...)

	if err != nil {
		return fmt.Errorf("failed create user %+v: %w", user, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed create user %+v: %w", user, entity.ErrAlreadyExists)
	}
	return nil
}

func (sp UserPsql) ReadByID(ctx context.Context, uid entity.UserID) (*entity.User, error) {
	query, args, err := sp.pg.Builder().
		Select(users.id).
		From(users.table).
		Where(sql.Eq{
			users.id:      uid,
			users.deleted: false,
		}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	row := sp.pg.Conn(ctx).QueryRow(ctx, query, args...)

	user := &entity.User{}
	if err := row.Scan(&user.ID); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, fmt.Errorf("user with id %v not found: %w", uid, entity.ErrNotFound)
		default:
			return nil, fmt.Errorf("failed read user by id %v: %w", uid, err)
		}
	}
	return user, nil
}
func (sp UserPsql) ReadAll(ctx context.Context) ([]*entity.User, error) {
	query, args, err := sp.pg.Builder().
		Select(users.id).
		From(users.table).
		Where(sql.Eq{
			users.deleted: false,
		}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := sp.pg.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed read all users: %w", err)
	}

	user := []*entity.User{}
	for rows.Next() {
		segment := &entity.User{}
		if err := rows.Scan(&segment.ID); err != nil {
			return nil, fmt.Errorf("failed scan user row: %w", err)
		}
		user = append(user, segment)
	}
	return user, nil
}

func (sp UserPsql) Delete(ctx context.Context, uid entity.UserID) error {
	query, args, err := sp.pg.Builder().
		Update(users.table).
		Set(users.deleted, true).
		Where(sql.Eq{
			users.id:      uid,
			users.deleted: false,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	tag, err := sp.pg.Conn(ctx).Exec(ctx, query, args...)

	if err != nil {
		return fmt.Errorf("failed deleting segment %v: %w", uid, err)
	}
	if tag.RowsAffected() == 0 {
		return entity.ErrNotFound
	}
	return nil
}

func (sp UserPsql) Prune(ctx context.Context) error {
	query, args, err := sp.pg.Builder().
		Delete(users.table).
		Where(sql.Eq{
			users.deleted: true,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	tag, err := sp.pg.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed prune users: %w", err)
	}
	log.Info().Msgf("prune %v users by %v", tag.RowsAffected(), tag.String())
	return nil
}
