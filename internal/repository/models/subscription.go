package models

import (
	"time"

	"github.com/scmbr/subscription-aggregator/internal/domain"
)

type Subscription struct {
	Id          string
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
}
type SubscriptionUpdate struct {
	ServiceName *string
	Price       *int
	UserID      *string
	StartDate   *time.Time
}
type GetAllSubscriptionsFilter struct {
	Limit       *int
	Offset      *int
	UserID      *string
	ServiceName *string
}

func SubscriptionDomainToModel(d *domain.Subscription) *Subscription {
	return &Subscription{
		Id:          d.Id,
		ServiceName: d.ServiceName,
		Price:       d.Price,
		UserID:      d.UserID,
		StartDate:   d.StartDate,
	}
}
func SubscriptionModelToDomain(m *Subscription) *domain.Subscription {
	return &domain.Subscription{
		Id:          m.Id,
		ServiceName: m.ServiceName,
		Price:       m.Price,
		UserID:      m.UserID,
		StartDate:   m.StartDate,
	}
}
