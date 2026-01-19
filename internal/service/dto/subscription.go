package dto

import "time"

type CreateSubscriptionInput struct {
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
	EndDate     *time.Time
}
type GetAllSubscriptionsInput struct {
	Limit  int
	Offset int
}
type GetAllSubscriptionsOutput struct {
	Total         int
	Subscriptions []GetSubscriptionOutput
}
type GetSubscriptionOutput struct {
	ID          int
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
	EndDate     *time.Time
}
type UpdateSubscriptionInput struct {
	ServiceName *string
	Price       *int
	UserID      *string
	StartDate   *time.Time
	EndDate     *time.Time
}
type GetTotalPriceInput struct {
	UserID      *string
	ServiceName *string
	StartDate   *time.Time
	EndDate     *time.Time
}
