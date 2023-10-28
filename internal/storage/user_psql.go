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
			<-time.After(60 * time.Second)
			log.Info().Err(users.Prune(ctx)).Msg("Pruning users")
		}
	}()
	return users, nil
}

func (sps UserPsql) Create(ctx context.Context, user *entity.User) error {
	query := `
    INSERT INTO users (id) VALUES ($1)
    ON CONFLICT (id) 
    DO UPDATE 
        SET deleted=FALSE
        WHERE users.deleted=TRUE`
	tag, err := sps.pg.Conn(ctx).Exec(ctx, query, user.ID)
	if err != nil {
		return fmt.Errorf("failed create user %+v: %w", user, err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("failed create user %+v: %w", user, ErrAlreadyExists)
	}
	return nil
}

func (sps UserPsql) ReadByID(ctx context.Context, uid entity.UserID) (*entity.User, error) {
	row := sps.pg.Conn(ctx).QueryRow(ctx, "SELECT id FROM users WHERE id=$1 AND deleted=FALSE", uid)

	user := &entity.User{}
	if err := row.Scan(&user.ID); err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, fmt.Errorf("user with id %v not found: %w", uid, ErrNotFound)
		default:
			return nil, fmt.Errorf("failed read user by id %v: %w", uid, err)
		}
	}
	return user, nil
}
func (sps UserPsql) ReadAll(ctx context.Context) ([]*entity.User, error) {
	rows, err := sps.pg.Conn(ctx).Query(ctx, "SELECT id FROM users WHERE deleted=FALSE")
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

func (sps UserPsql) Delete(ctx context.Context, uid entity.UserID) error {
	tag, err := sps.pg.Conn(ctx).Exec(ctx, "UPDATE users SET deleted=TRUE WHERE id=$1 AND deleted=FALSE", uid)
	if err != nil {
		return fmt.Errorf("failed deleting segment %v: %w", uid, err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (sps UserPsql) Prune(ctx context.Context) error {
	if _, err := sps.pg.Conn(ctx).Exec(ctx, "DELETE FROM users WHERE deleted=TRUE"); err != nil {
		return fmt.Errorf("failed prune users: %w", err)
	}
	return nil
}
