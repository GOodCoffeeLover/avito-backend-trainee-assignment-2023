package main

import (
	"context"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/app/config"
	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/pkg/postgres"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	pg, trm, err := postgres.New(ctx, config.New().ConnString)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect postgres")
	}
	defer pg.Close(ctx)
	m := newMigrator(ctx, pg, zerolog.TraceLevel)
	err = trm.Do(ctx, func(ctx context.Context) error {
		if err := m.dropAssigmentsTable(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to drop users table")
			return err
		}
		if err := m.dropEventsTable(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to drop events table")
			return err
		}
		if err := m.dropSegmentsTable(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to drop segments table")
			return err
		}
		if err := m.dropUsersTable(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to drop users table")
			return err
		}

		if err := m.createSegmentsTable(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to create segments table")
			return err
		}
		if err := m.createUsersTable(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to create users table")
			return err
		}
		if err := m.createAssigmentsTable(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to create users table")
			return err
		}
		if err := m.createEventsTable(ctx); err != nil {
			log.Fatal().Err(err).Msg("failed to create events table")
			return err
		}
		return nil
	})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to exec migration")
	}
}

func newMigrator(ctx context.Context, pg *postgres.Postgres, level zerolog.Level) *migrator {

	log := log.Logger.With().Caller().Timestamp().Str("component", "migration").Logger()
	return &migrator{
		pg:  pg,
		log: &log,
	}
}

type migrator struct {
	pg  *postgres.Postgres
	log *zerolog.Logger
}

func (m migrator) createSegmentsTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS segments (
        name varchar(40) PRIMARY KEY UNIQUE,
        deleted bool DEFAULT FALSE
    )`

	tag, err := m.pg.Conn(ctx).Exec(ctx, query)
	m.log.Info().Err(err).Msg(tag.String())
	return err
}
func (m migrator) dropSegmentsTable(ctx context.Context) error {
	query := `DROP TABLE IF EXISTS segments`

	tag, err := m.pg.Conn(ctx).Exec(ctx, query)
	m.log.Info().Err(err).Msg(tag.String())
	return err
}

func (m migrator) createUsersTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS users (
        id integer PRIMARY KEY UNIQUE CONSTRAINT uint CHECK (id > 0),
        deleted bool DEFAULT FALSE
    )`

	tag, err := m.pg.Conn(ctx).Exec(ctx, query)
	m.log.Info().Err(err).Msg(tag.String())
	return err
}
func (m migrator) dropUsersTable(ctx context.Context) error {
	query := `DROP TABLE IF EXISTS users`

	tag, err := m.pg.Conn(ctx).Exec(ctx, query)
	m.log.Info().Err(err).Msg(tag.String())
	return err
}

func (m migrator) createAssigmentsTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS assignments (
        user_id integer REFERENCES users(id) ON DELETE CASCADE,
        segment_name VARCHAR(40) REFERENCES segments(name) ON DELETE CASCADE,
        deleted bool DEFAULT FALSE, 
        PRIMARY KEY (user_id,segment_name)
    )`

	tag, err := m.pg.Conn(ctx).Exec(ctx, query)
	m.log.Info().Err(err).Msg(tag.String())
	return err
}
func (m migrator) dropAssigmentsTable(ctx context.Context) error {
	query := `DROP TABLE IF EXISTS assignments`

	tag, err := m.pg.Conn(ctx).Exec(ctx, query)
	m.log.Info().Err(err).Msg(tag.String())
	return err
}

func (m migrator) createEventsTable(ctx context.Context) error {
	query := `CREATE TABLE IF NOT EXISTS events (
        user_id integer CONSTRAINT uint CHECK (user_id > 0),
        segment VARCHAR(40),
		operation VARCHAR(40),
		timestamp timestamp
		)`

	tag, err := m.pg.Conn(ctx).Exec(ctx, query)
	m.log.Info().Err(err).Msg(tag.String())
	return err
}
func (m migrator) dropEventsTable(ctx context.Context) error {
	query := `DROP TABLE IF EXISTS events`

	tag, err := m.pg.Conn(ctx).Exec(ctx, query)
	m.log.Info().Err(err).Msg(tag.String())
	return err
}
