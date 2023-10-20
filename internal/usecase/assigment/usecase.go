package assigment

import (
	"context"
	"fmt"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
)

type UseCase struct {
	segments   SegmentStorage
	users      UserStorage
	assigments AssignmentStorage
	trxManager *manager.Manager // for history transactions
}

func New(segments SegmentStorage, users UserStorage, assigments AssignmentStorage, trxManager *manager.Manager) *UseCase {
	return &UseCase{
		segments:   segments,
		users:      users,
		assigments: assigments,
		trxManager: trxManager,
	}
}

func (uc *UseCase) ReadByUserID(ctx context.Context, uid entity.UserID) (*entity.Assignment, error) {

	err := uc.trxManager.Do(ctx, func(ctx context.Context) error {
		if _, err := uc.users.ReadByID(ctx, uid); err != nil {
			return fmt.Errorf("user with id %v %w", uid, entity.ErrNotFound)
		}

		return nil
	})
	return nil, err
}

func AssignSegments(ctx context.Context, uid entity.UserID, segments []entity.SegmentName) error {
	return nil
}

func UnassignSegments(ctx context.Context, uid entity.UserID, segments []entity.SegmentName) error {
	return nil
}
