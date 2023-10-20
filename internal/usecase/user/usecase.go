package user

import (
	"context"
	"fmt"

	"github.com/GOodCoffeeLover/avito-backend-trainee-assignment-2023/internal/entity"
	"github.com/avito-tech/go-transaction-manager/trm/manager"
)

type UseCase struct {
	users      UserStorage
	trxManager *manager.Manager
}

func New(users UserStorage, trxManager *manager.Manager) *UseCase {
	return &UseCase{
		users:      users,
		trxManager: trxManager,
	}
}

func (uc *UseCase) Create(ctx context.Context, uid entity.UserID) error {
	user, err := entity.NewUser(uid)
	if err != nil {
		return fmt.Errorf("failed to build new user: %w", err)
	}
	err = uc.trxManager.Do(ctx, func(ctx context.Context) error {
		return uc.users.Create(ctx, user)
	})
	if err != nil {
		return fmt.Errorf("failed to create user %v: %w", uid, err)
	}
	return nil
}

func (uc *UseCase) Read(ctx context.Context, uid entity.UserID) (*entity.User, error) {
	var user *entity.User
	err := uc.trxManager.Do(ctx, func(ctx context.Context) error {
		var err error
		user, err = uc.users.ReadByID(ctx, uid)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read user %v: %w", uid, err)
	}
	return user, nil
}

func (uc *UseCase) ReadAll(ctx context.Context) ([]*entity.User, error) {
	var segs []*entity.User
	err := uc.trxManager.Do(ctx, func(ctx context.Context) error {
		var err error
		segs, err = uc.users.ReadAll(ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("failed to read all users: %w", err)
	}
	return segs, nil
}

func (uc *UseCase) Delete(ctx context.Context, uid entity.UserID) error {
	err := uc.users.Delete(ctx, uid)
	if err != nil {
		return fmt.Errorf("failed to delete %v: %w", uid, err)
	}
	return nil
}
