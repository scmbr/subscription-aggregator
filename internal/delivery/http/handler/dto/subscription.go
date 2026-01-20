package dto

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int        `json:"price" binding:"required,gte=0"`
	UserID      string     `json:"user_id" binding:"required,uuid4"`
	StartDate   MonthYear  `json:"start_date" binding:"required"`
	EndDate     *MonthYear `json:"end_date" binding:"omitempty"`
}
type CreateSubscriptionResponse struct {
	Id string `json:"subscription_id"`
}
type GetAllSubscriptionsResponse struct {
	Total         int                       `json:"total"`
	Subscriptions []GetSubscriptionResponse `json:"subscriptions"`
}
type GetSubscriptionResponse struct {
	Id          string     `json:"subscription_id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	UserID      string     `json:"user_id"`
	StartDate   MonthYear  `json:"start_date"`
	EndDate     *MonthYear `json:"end_date"`
}
type UpdateSubscriptionRequest struct {
	ServiceName *string    `json:"service_name"  binding:"omitempty"`
	Price       *int       `json:"price"  binding:"omitempty,gte=0"`
	UserID      *string    `json:"user_id"  binding:"omitempty"`
	StartDate   *MonthYear `json:"start_date"  binding:"omitempty"`
	EndDate     *MonthYear `json:"end_date"  binding:"omitempty"`
}

type GetTotalPriceRequest struct {
	UserID      *string    `form:"user_id"`
	ServiceName *string    `form:"service_name"`
	StartDate   *MonthYear `form:"start_date" binding:"required"`
	EndDate     *MonthYear `form:"end_date" binding:"required"`
}
type GetTotalPriceResponse struct {
	TotalPrice int `json:"total_price"`
}
