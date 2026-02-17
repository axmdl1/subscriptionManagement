package service

import (
	"context"

	"subscriptionManagement/internal/model"
)

type SubscriptionService interface {
	Create(ctx context.Context, sub model.Subscription) (*model.Subscription, error)
	GetByID(ctx context.Context, id string) (*model.Subscription, error)
	List(ctx context.Context) ([]model.Subscription, error)
	Update(ctx context.Context, sub model.Subscription) (*model.Subscription, error)
	Delete(ctx context.Context, id string) error

	CalculateTotal(ctx context.Context, filter model.SubscriptionFilter) (int, error)
}
