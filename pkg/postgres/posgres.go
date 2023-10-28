package postgres

import (
	"context"
	"fmt"

	sql "github.com/Masterminds/squirrel"
	"github.com/avito-tech/go-transaction-manager/pgxv5"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(ctx context.Context, connString string) (*Postgres, *manager.Manager, error) {

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to psql: %w", err)
	}
	trm, err := manager.New(pgxv5.NewDefaultFactory(pool))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create transaction manager: %w", err)
	}
	builder := sql.StatementBuilder.PlaceholderFormat(sql.Dollar)

	return &Postgres{
		pgxPool: pool,
		getter:  pgxv5.DefaultCtxGetter,
		builder: builder,
	}, trm, nil

}

type Postgres struct {
	pgxPool *pgxpool.Pool
	getter  *pgxv5.CtxGetter
	builder sql.StatementBuilderType
}

func (p Postgres) Builder() sql.StatementBuilderType {
	return p.builder
}
func (p Postgres) Conn(ctx context.Context) pgxv5.Tr {
	return p.getter.DefaultTrOrDB(ctx, p.pgxPool)
}

func (p Postgres) Close(ctx context.Context) {
	p.pgxPool.Close()
}
