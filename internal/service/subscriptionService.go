package service

import (
	"context"
	"errors"
	"time"

	"subscriptionManagement/internal/model"
	"subscriptionManagement/internal/repository"
	"subscriptionManagement/internal/utils"
)

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(r repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: r}
}

func (s *subscriptionService) Create(ctx context.Context, sub model.Subscription) (*model.Subscription, error) {

	if sub.Price < 0 {
		return nil, errors.New("price must be positive")
	}

	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		return nil, errors.New("end_date before start_date")
	}

	err := s.repo.Create(ctx, &sub)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *subscriptionService) GetByID(ctx context.Context, id string) (*model.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) List(ctx context.Context) ([]model.Subscription, error) {
	return s.repo.List(ctx)
}

func (s *subscriptionService) Update(ctx context.Context, sub model.Subscription) (*model.Subscription, error) {
	if sub.EndDate != nil && sub.EndDate.Before(sub.StartDate) {
		return nil, errors.New("end_date before start_date")
	}

	err := s.repo.Update(ctx, &sub)
	if err != nil {
		return nil, err
	}

	return &sub, nil
}

func (s *subscriptionService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *subscriptionService) CalculateTotal(
	ctx context.Context,
	filter model.SubscriptionFilter,
) (int, error) {

	subs, err := s.repo.GetAllForTotal(ctx, filter)
	if err != nil {
		return 0, err
	}

	total := 0

	for _, sub := range subs {

		subEnd := time.Now()
		if sub.EndDate != nil {
			subEnd = *sub.EndDate
		}

		overlapStart, overlapEnd, ok := utils.OverlapPeriod(
			sub.StartDate,
			subEnd,
			filter.From,
			filter.To,
		)

		if !ok {
			continue
		}

		months := utils.MonthsBetween(overlapStart, overlapEnd)
		total += months * sub.Price
	}

	return total, nil
}
