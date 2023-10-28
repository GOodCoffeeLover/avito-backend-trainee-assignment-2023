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

func (uc *UseCase) ReadByUserID(ctx context.Context, uid entity.UserID) ([]*entity.Assignment, error) {
	var assigments []*entity.Assignment
	err := uc.trxManager.Do(ctx, func(ctx context.Context) error {
		var err error
		if _, err = uc.users.ReadByID(ctx, uid); err != nil {
			return fmt.Errorf("unknown user with id %v %w", uid, entity.ErrNotFound)
		}
		assigments, err = uc.assigments.ReadByUserID(ctx, uid)
		if err != nil {
			return fmt.Errorf("failed get assigments for userid(%v) %w", uid, entity.ErrNotFound)
		}
		return nil
	})
	return assigments, err
}

func (uc *UseCase) SetToUserByID(ctx context.Context, uid entity.UserID, segments []entity.SegmentName) error {
	return uc.trxManager.Do(ctx, func(ctx context.Context) error {
		var err error
		if _, err = uc.users.ReadByID(ctx, uid); err != nil {
			return fmt.Errorf("unknown user with id %v %w", uid, entity.ErrNotFound)
		}
		for _, segment := range segments {
			if _, err = uc.segments.ReadByName(ctx, segment); err != nil {
				return fmt.Errorf("unknown segment %v %w", segment, entity.ErrNotFound)
			}
		}

		for _, segment := range segments {
			assigment, err := entity.NewAssignment(uid, segment)
			if err != nil {
				return fmt.Errorf("failed to create assigment to segment %v : %w", segment, err)
			}
			if err = uc.assigments.Create(ctx, assigment); err != nil {
				return fmt.Errorf("failed to save assigment %v: %w", assigment, err)
			}
		}
		return nil
	})
}

func (uc *UseCase) UnsetToUserByID(ctx context.Context, uid entity.UserID, segments []entity.SegmentName) error {
	return uc.trxManager.Do(ctx, func(ctx context.Context) error {
		assigments, err := uc.assigments.ReadByUserID(ctx, uid)
		if err != nil {
			return fmt.Errorf("failed to get assigments for user %v: %w", uid, err)
		}
		segmentsForDeletion := map[entity.SegmentName]struct{}{}
		for _, segment := range segments {
			segmentsForDeletion[segment] = struct{}{}

		}
		for _, assigment := range assigments {
			if _, ok := segmentsForDeletion[assigment.Segment]; ok {
				uc.assigments.Delete(ctx, assigment)
				delete(segmentsForDeletion, assigment.Segment)
			}
		}
		if len(segmentsForDeletion) != 0 {
			segs := make([]entity.SegmentName, 0, len(segmentsForDeletion))
			for segment := range segmentsForDeletion {
				segs = append(segs, segment)
			}
			return fmt.Errorf("unassigned segs: %v", segs)
		}
		return nil
	})
}
