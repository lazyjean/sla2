package repository

import (
	"context"

	"github.com/lazyjean/sla2/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
	FindByAccount(ctx context.Context, account string) (*entity.User, error)
}
