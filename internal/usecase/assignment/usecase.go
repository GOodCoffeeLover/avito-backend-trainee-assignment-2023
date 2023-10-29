package assignment

import (
	"context"
	"fmt"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
)

type UseCase struct {
	segments    SegmentStorage
	users       UserStorage
	assignments AssignmentStorage
	trxManager  *manager.Manager // for history transactions
}

func New(segments SegmentStorage, users UserStorage, assignments AssignmentStorage, trxManager *manager.Manager) *UseCase {
	return &UseCase{
		segments:    segments,
		users:       users,
		assignments: assignments,
		trxManager:  trxManager,
	}
}

func (uc *UseCase) ReadByUserID(ctx context.Context, uid entity.UserID) ([]*entity.Assignment, error) {
	var assignments []*entity.Assignment
	err := uc.trxManager.Do(ctx, func(ctx context.Context) error {
		var err error
		if _, err = uc.users.ReadByID(ctx, uid); err != nil {
			return fmt.Errorf("unknown user with id %v %w", uid, entity.ErrNotFound)
		}
		assignments, err = uc.assignments.ReadByUserID(ctx, uid)
		if err != nil {
			return fmt.Errorf("failed get assignments for userid(%v): %w: %w:", uid, err, entity.ErrNotFound)
		}
		return nil
	})
	return assignments, err
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
			assignment, err := entity.NewAssignment(uid, segment)
			if err != nil {
				return fmt.Errorf("failed to create assignment to segment %v : %w", segment, err)
			}
			if err = uc.assignments.Create(ctx, assignment); err != nil {
				return fmt.Errorf("failed to save assignment %v: %w", assignment, err)
			}
		}
		return nil
	})
}

func (uc *UseCase) UnsetToUserByID(ctx context.Context, uid entity.UserID, segments []entity.SegmentName) error {
	return uc.trxManager.Do(ctx, func(ctx context.Context) error {
		assignments, err := uc.assignments.ReadByUserID(ctx, uid)
		if err != nil {
			return fmt.Errorf("failed to get assignments for user %v: %w", uid, err)
		}
		segmentsForDeletion := map[entity.SegmentName]struct{}{}
		for _, segment := range segments {
			segmentsForDeletion[segment] = struct{}{}

		}
		for _, assignment := range assignments {
			if _, ok := segmentsForDeletion[assignment.Segment]; ok {
				uc.assignments.Delete(ctx, assignment)
				delete(segmentsForDeletion, assignment.Segment)
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
