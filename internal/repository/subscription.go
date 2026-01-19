package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/scmbr/subscription-aggregator/internal/domain"
	"github.com/scmbr/subscription-aggregator/internal/repository/models"
)

type SubscriptionRepo struct {
	db *sqlx.DB
}

func NewSubscriptionRepository(db *sqlx.DB) *SubscriptionRepo {
	return &SubscriptionRepo{
		db: db,
	}
}
func (r *SubscriptionRepo) Create(ctx context.Context, input *domain.Subscription) error {
	_, err := r.db.ExecContext(ctx, `
    INSERT INTO subscriptions (id, service_name, price, user_id, start_date, end_date)
    VALUES ($1, $2, $3, $4, $5, $6)
`, input.Id, input.ServiceName, input.Price, input.UserID, input.StartDate, input.EndDate)
	if err != nil {
		return fmt.Errorf("subscriptionRepo.Create:%w", err)
	}
	return nil
}
func (r *SubscriptionRepo) GetAll(ctx context.Context, limit, offset int) ([]*domain.Subscription, int, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date
              FROM subscriptions`

	args := []interface{}{}

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}
	if offset > 0 {
		query += " OFFSET ?"
		args = append(args, offset)
	}

	countQuery := "SELECT COUNT(*) FROM subscriptions"

	query = sqlx.Rebind(sqlx.DOLLAR, query)
	countQuery = sqlx.Rebind(sqlx.DOLLAR, countQuery)

	subscriptions := make([]*models.Subscription, 0)
	if err := r.db.SelectContext(ctx, &subscriptions, query, args...); err != nil {
		return nil, 0, fmt.Errorf("subscriptionRepo.GetAll: %w", err)
	}

	var count int
	if err := r.db.GetContext(ctx, &count, countQuery); err != nil {
		return nil, 0, fmt.Errorf("subscriptionRepo.GetAll: %w", err)
	}

	subscriptionsDomain := make([]*domain.Subscription, 0, len(subscriptions))
	for _, s := range subscriptions {
		subscriptionsDomain = append(subscriptionsDomain, models.SubscriptionModelToDomain(s))
	}

	return subscriptionsDomain, count, nil
}

func (r *SubscriptionRepo) GetById(ctx context.Context, id string) (*domain.Subscription, error) {
	var subscription models.Subscription
	query := "SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1"

	if err := r.db.GetContext(ctx, &subscription, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return models.SubscriptionModelToDomain(&subscription), nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, id string, input models.SubscriptionUpdate) error {
	set := []string{}
	args := []interface{}{}
	idx := 1

	if input.ServiceName != nil {
		set = append(set, fmt.Sprintf("service_name = $%d", idx))
		args = append(args, *input.ServiceName)
		idx++
	}
	if input.Price != nil {
		set = append(set, fmt.Sprintf("price = $%d", idx))
		args = append(args, *input.Price)
		idx++
	}
	if input.StartDate != nil {
		set = append(set, fmt.Sprintf("start_date = $%d", idx))
		args = append(args, *input.StartDate)
		idx++
	}

	if len(set) == 0 {
		return nil
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE subscriptions SET %s WHERE id = $%d", strings.Join(set, ", "), idx)

	res, err := r.db.ExecContext(ctx, query, args...)
	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrNotFound
	}
	return err
}

func (r *SubscriptionRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, "DELETE FROM subscriptions WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("subscriptionRepo.Delete:%w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrNotFound
	}

	return nil
}
func (r *SubscriptionRepo) GetTotalPrice(ctx context.Context, filter models.GetTotalPriceFilter) (int, error) {
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions`
	where := []string{}
	args := []interface{}{}

	if filter.UserID != nil {
		where = append(where, "user_id = ?")
		args = append(args, *filter.UserID)
	}
	if filter.ServiceName != nil {
		where = append(where, "service_name = ?")
		args = append(args, *filter.ServiceName)
	}
	if filter.StartDate != nil {
		where = append(where, "start_date <= ?")
		args = append(args, *filter.StartDate)
	}
	if filter.EndDate != nil {
		where = append(where, "end_date >= ?")
		args = append(args, *filter.EndDate)
	}

	if len(where) > 0 {
		query += " WHERE " + strings.Join(where, " AND ")
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	var total int
	if err := r.db.GetContext(ctx, &total, query, args...); err != nil {
		return 0, fmt.Errorf("subscriptionRepo.GetTotalPrice: %w", err)
	}

	return total, nil
}
