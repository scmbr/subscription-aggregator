package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/scmbr/subscription-aggregator/internal/domain"
	"github.com/scmbr/subscription-aggregator/internal/repository/models"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, input *domain.Subscription) error
	GetAll(ctx context.Context, filter models.GetAllSubscriptionsFilter) ([]*domain.Subscription, int, error)
	GetById(ctx context.Context, id string) (*domain.Subscription, error)
	Update(ctx context.Context, id string, input models.SubscriptionUpdate) error
	Delete(ctx context.Context, id string) error
}
type Repository struct {
	Subscription SubscriptionRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Subscription: NewSubscriptionRepository(db),
	}
}
