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

type EventPsql struct {
	pg *postgres.Postgres
}

func NewEventPsql(ctx context.Context, pg *postgres.Postgres) (*EventPsql, error) {
	if pg == nil {
		return nil, fmt.Errorf("nil potgresql client")
	}
	events := &EventPsql{
		pg: pg,
	}
	go func() {
		ctx := context.Background()
		for {
			<-time.After(120 * time.Second)
			log.Info().Err(events.Prune(ctx)).Msg("Pruned events")
		}
	}()
	return events, nil
}

func (ep EventPsql) Create(ctx context.Context, event *entity.Event) error {
	query, args, err := ep.pg.Builder().
		Insert(events.table).
		Columns(events.user, events.segment, events.operation, events.timestamp).
		Values(event.User, event.Segment, event.Op, event.Ts).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = ep.pg.Conn(ctx).Exec(ctx, query, args...)

	if err != nil {
		return fmt.Errorf("failed to exec query %v %v : %w", query, args, err)
	}
	return nil
}

func (ep EventPsql) ReadByUserID(ctx context.Context, uid entity.UserID, start, end *time.Time) ([]*entity.Event, error) {
	query, args, err := ep.pg.Builder().
		Select(events.user, events.segment, events.operation, events.timestamp).
		From(events.table).
		Where(sql.Or{
			sql.Eq{events.user: uid},
			sql.Eq{events.user: nil},
		}).
		Where(sql.GtOrEq{events.timestamp: start}).
		Where(sql.LtOrEq{events.timestamp: end}).
		LeftJoin("( SELECT id FROM users WHERE deleted = FALSE ) AS users ON events.user_id = users.id").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}
	rows, err := ep.pg.Conn(ctx).Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to exec query %v %v : %w", query, args, err)
	}

	userEvents := []*entity.Event{}
	for rows.Next() {
		event := entity.Event{}
		err = rows.Scan(&event.User, &event.Segment, &event.Op, &event.Ts)
		if err != nil {
			return nil, fmt.Errorf("failed to scan events: %w", err)
		}
		userEvents = append(userEvents, &event)
	}
	return userEvents, nil
}

func (ep EventPsql) Prune(ctx context.Context) error {
	query, args, err := ep.pg.Builder().
		Delete(events.table).
		Where(sql.Lt{
			events.timestamp: true,
		}).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	tag, err := ep.pg.Conn(ctx).Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed prune events: %w", err)
	}
	log.Info().Msgf("prune %v events by %v", tag.RowsAffected(), tag.String())
	return nil
}
