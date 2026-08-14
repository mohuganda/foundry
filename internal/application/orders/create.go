package orders

import (
	"context"
	domain "github.com/jabahum/foundry/internal/domain/orders"
)

type CreateInput struct{}

type CreateUseCase struct { repo domain.Repository }

func NewCreateUseCase(repo domain.Repository) *CreateUseCase { return &CreateUseCase{repo: repo} }

func (u *CreateUseCase) Execute(ctx context.Context, input CreateInput) (*domain.Orders, error) {
	return nil, nil
}
