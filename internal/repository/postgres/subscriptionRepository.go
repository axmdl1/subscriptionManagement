package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"subscriptionManagement/internal/model"
)

type SubscriptionRepository struct {
	pool *pgxpool.Pool
}

func NewSubscriptionRepository(pool *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{pool: pool}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *model.Subscription) error {
	query := `
	INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`

	return r.pool.QueryRow(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&sub.ID)
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id string) (*model.Subscription, error) {
	query := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	WHERE id = $1
	`

	var sub model.Subscription

	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)

	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (r *SubscriptionRepository) List(ctx context.Context) ([]model.Subscription, error) {
	query := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	ORDER BY created_at DESC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Subscription

	for rows.Next() {
		var sub model.Subscription

		if err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		); err != nil {
			return nil, err
		}

		result = append(result, sub)
	}

	return result, rows.Err()
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *model.Subscription) error {
	query := `
	UPDATE subscriptions
	SET service_name = $1,
	    price = $2,
	    user_id = $3,
	    start_date = $4,
	    end_date = $5,
	    updated_at = NOW()
	WHERE id = $6
	`

	cmd, err := r.pool.Exec(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id string) error {
	cmd, err := r.pool.Exec(ctx,
		`DELETE FROM subscriptions WHERE id = $1`, id)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *SubscriptionRepository) GetAllForTotal(
	ctx context.Context,
	filter model.SubscriptionFilter,
) ([]model.Subscription, error) {

	query := `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	WHERE start_date <= $1
	AND (end_date IS NULL OR end_date >= $2)
	`

	args := []any{filter.To, filter.From}
	argPos := 3

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *filter.UserID)
		argPos++
	}

	if filter.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argPos)
		args = append(args, *filter.ServiceName)
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.Subscription

	for rows.Next() {
		var sub model.Subscription
		if err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		); err != nil {
			return nil, err
		}
		result = append(result, sub)
	}

	return result, rows.Err()
}
