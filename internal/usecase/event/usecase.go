package event

import (
	"context"
	"fmt"
	"time"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type UseCase struct {
	events     EventStorage
	users      UserStorage
	trxManager *manager.Manager // TODO: move to infra layer
	log        zerolog.Logger
}

func New(users UserStorage, events EventStorage, trxManager *manager.Manager) *UseCase {
	log := log.Logger.With().Str("level", "usecase").Str("component", "assigments").Logger()
	return &UseCase{
		events:     events,
		users:      users,
		trxManager: trxManager,
		log:        log,
	}
}

type InReadEventsByUserID struct {
	UID   entity.UserID
	Start *time.Time
	End   *time.Time
}

func (uc *UseCase) ReadByUserID(ctx context.Context, in *InReadEventsByUserID) ([]*entity.Event, error) {
	if in.Start.After(*in.End) {
		return nil, fmt.Errorf("%w: start %v after end %v", entity.ErrInvalidArgument, in.Start, in.End)
	}
	var events []*entity.Event
	err := uc.trxManager.Do(ctx, func(ctx context.Context) error {
		var err error
		if _, err = uc.users.ReadByID(ctx, in.UID); err != nil {
			return fmt.Errorf("unknown user with id %v %w", in.UID, entity.ErrNotFound)
		}
		events, err = uc.events.ReadByUserID(ctx, in.UID, in.Start, in.End)
		if err != nil {
			return fmt.Errorf("failed get assignments for userid(%v): %w: %w", in.UID, err, entity.ErrNotFound)
		}
		return nil
	})
	return events, err
}
