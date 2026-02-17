package repository

import (
	"context"

	"subscriptionManagement/internal/model"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *model.Subscription) error
	GetByID(ctx context.Context, id string) (*model.Subscription, error)
	List(ctx context.Context) ([]model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) error
	Delete(ctx context.Context, id string) error
	GetAllForTotal(ctx context.Context, filter model.SubscriptionFilter) ([]model.Subscription, error)
}
