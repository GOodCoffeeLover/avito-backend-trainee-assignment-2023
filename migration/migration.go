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
	pg, err := postgres.New(ctx, config.New().ConnString)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect postgres")
	}
	defer pg.Close(ctx)
	m := newMigrator(ctx, pg, zerolog.TraceLevel)
	if err = m.dropSegmentsTable(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to drop segments table")
	}
	if err = m.createSegmentsTable(ctx); err != nil {
		log.Fatal().Err(err).Msg("failed to create segments table")
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

	tag, err := m.pg.Conn().Exec(ctx, query)
	m.log.Info().Err(err).Msg(tag.String())
	return err
}
func (m migrator) dropSegmentsTable(ctx context.Context) error {
	query := `DROP TABLE IF EXISTS segments`

	tag, err := m.pg.Conn().Exec(ctx, query)
	m.log.Info().Err(err).Msg(tag.String())
	return err
}
