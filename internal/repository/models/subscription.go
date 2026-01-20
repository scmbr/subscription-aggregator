package models

import (
	"time"

	"github.com/scmbr/subscription-aggregator/internal/domain"
)

type Subscription struct {
	Id          string     `db:"id"`
	ServiceName string     `db:"service_name"`
	Price       int        `db:"price"`
	UserID      string     `db:"user_id"`
	StartDate   time.Time  `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
}

type SubscriptionUpdate struct {
	ServiceName *string    `db:"service_name"`
	Price       *int       `db:"price"`
	UserID      *string    `db:"user_id"`
	StartDate   *time.Time `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
}

type GetTotalPriceFilter struct {
	UserID      *string    `db:"user_id"`
	ServiceName *string    `db:"service_name"`
	StartDate   *time.Time `db:"start_date"`
	EndDate     *time.Time `db:"end_date"`
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
