package domain

import "time"

type Subscription struct {
	Id          string
	ServiceName string
	Price       int
	UserID      string
	StartDate   time.Time
}
