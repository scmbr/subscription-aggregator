package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/scmbr/subscription-aggregator/internal/domain"
	"github.com/scmbr/subscription-aggregator/internal/repository"
	"github.com/scmbr/subscription-aggregator/internal/repository/models"
	"github.com/scmbr/subscription-aggregator/internal/service/dto"
)

type SubscriptionSvc struct {
	subscriptionRepo repository.SubscriptionRepository
}

func NewSubscriptionService(subscriptionRepo repository.SubscriptionRepository) *SubscriptionSvc {
	return &SubscriptionSvc{
		subscriptionRepo: subscriptionRepo,
	}
}
func (s *SubscriptionSvc) CreateSubscription(ctx context.Context, input *dto.CreateSubscriptionInput) (string, error) {
	id := uuid.NewString()
	subscriptionDomain, err := domain.NewSubscription(
		id,
		input.ServiceName,
		input.Price,
		input.UserID,
		input.StartDate,
		input.EndDate,
	)
	if err != nil {
		return "", err
	}
	err = s.subscriptionRepo.Create(ctx, subscriptionDomain)
	if err != nil {
		return "", err
	}
	return id, nil
}
func (s *SubscriptionSvc) GetAllSubscriptions(ctx context.Context, input dto.GetAllSubscriptionsInput) (*dto.GetAllSubscriptionsOutput, error) {
	subscriptions, total, err := s.subscriptionRepo.GetAll(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, err
	}
	subscriptionsDTO := make([]*dto.GetSubscriptionOutput, 0, len(subscriptions))
	for _, s := range subscriptions {
		subscriptionsDTO = append(subscriptionsDTO, &dto.GetSubscriptionOutput{
			ID:          s.Id,
			ServiceName: s.ServiceName,
			Price:       s.Price,
			UserID:      s.UserID,
			StartDate:   s.StartDate,
			EndDate:     s.EndDate,
		})
	}
	return &dto.GetAllSubscriptionsOutput{
		Total:         total,
		Subscriptions: subscriptionsDTO,
	}, nil
}
func (s *SubscriptionSvc) GetSubscriptionById(ctx context.Context, id string) (*dto.GetSubscriptionOutput, error) {
	subscription, err := s.subscriptionRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrSubscriptionNotFound
		}
		return nil, err
	}

	return &dto.GetSubscriptionOutput{
		ID:          subscription.Id,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID,
		StartDate:   subscription.StartDate,
		EndDate:     subscription.EndDate,
	}, nil
}
func (s *SubscriptionSvc) UpdateSubscriptionById(ctx context.Context, id string, input *dto.UpdateSubscriptionInput) error {
	if err := s.subscriptionRepo.Update(ctx, id, models.SubscriptionUpdate{
		ServiceName: input.ServiceName,
		Price:       input.Price,
		UserID:      input.UserID,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	}); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		return err
	}
	return nil
}
func (s *SubscriptionSvc) DeleteSubscriptionById(ctx context.Context, id string) error {
	if err := s.subscriptionRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return ErrSubscriptionNotFound
		}
		return err
	}
	return nil
}
func (s *SubscriptionSvc) GetSubscriptionsTotalPrice(ctx context.Context, input *dto.GetTotalPriceInput) (int, error) {
	total, err := s.subscriptionRepo.GetTotalPrice(ctx, models.GetTotalPriceFilter{
		ServiceName: input.ServiceName,
		UserID:      input.UserID,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	})
	if err != nil {
		return 0, err
	}
	return total, nil
}
