package domain

import (
	"fmt"
	"time"
)

type Subscription struct {
	Id          string
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
	EndDate     *time.Time
}

func NewSubscription(id, serviceName string, price int, userID string, startDate time.Time, endDate *time.Time) (*Subscription, error) {
	if endDate != nil && endDate.Before(startDate) {
		return nil, fmt.Errorf("domain.NewSubscription invalid startDate and endDate")
	}
	if price < 0 {
		return nil, fmt.Errorf("domain.NewSubscription invalid price")
	}
	return &Subscription{
		Id:          id,
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}
