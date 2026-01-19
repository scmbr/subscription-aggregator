package service

import (
	"context"

	"github.com/scmbr/subscription-aggregator/internal/repository"
	"github.com/scmbr/subscription-aggregator/internal/service/dto"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, input *dto.CreateSubscriptionInput) error
	GetAllSubscriptions(ctx context.Context, input *dto.GetAllSubscriptionsInput) (*dto.GetAllSubscriptionsOutput, error)
	GetSubscriptionById(ctx context.Context, id string) (*dto.GetSubscriptionOutput, error)
	UpdateSubscriptionById(ctx context.Context, id string, input *dto.UpdateSubscriptionInput) error
	DeleteSubscriptionById(ctx context.Context, id string) error
	GetSubscriptionsTotalPrice(ctx context.Context, input *dto.GetTotalPriceInput) (int, error)
}
type Service struct {
	Subscription SubscriptionService
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Subscription: NewSubscriptionService(repo.Subscription),
	}
}
