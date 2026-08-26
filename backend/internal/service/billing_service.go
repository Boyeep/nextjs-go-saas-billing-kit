package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Boyeep/nextjs-go-saas-billing-kit/backend/internal/domain"
	"github.com/Boyeep/nextjs-go-saas-billing-kit/backend/internal/repository"
)

var ErrBillingForbidden = errors.New("billing access forbidden")

type BillingProvider interface {
	Checkout(context.Context, string, string, string) (string, error)
	Portal(context.Context, string, string) (string, error)
	ParseWebhook([]byte, string) (BillingEvent, error)
}

type BillingEvent struct {
	ID             string
	Type           string
	WorkspaceID    string
	CustomerID     string
	SubscriptionID string
	Plan           string
	Status         string
	PeriodEnd      *time.Time
	CancelAtEnd    bool
}

type BillingService struct {
	repository repository.BillingRepository
	provider   BillingProvider
}

func NewBillingService(repo repository.BillingRepository, provider BillingProvider) *BillingService {
	return &BillingService{repository: repo, provider: provider}
}

func (s *BillingService) Workspaces(ctx context.Context, userID string) ([]domain.Workspace, error) {
	return s.repository.ListWorkspaces(ctx, userID)
}

func (s *BillingService) Summary(ctx context.Context, workspaceID, userID string) (domain.BillingSummary, error) {
	summary, err := s.repository.GetSummary(ctx, workspaceID, userID)
	if errors.Is(err, repository.ErrBillingNotFound) {
		return summary, ErrBillingForbidden
	}
	if err != nil {
		return summary, err
	}
	summary.Entitlements = entitlements(summary.Subscription.Plan, summary.Subscription.Status, summary.Subscription.GracePeriodEnd)
	return summary, nil
}

func (s *BillingService) Checkout(ctx context.Context, workspaceID, userID, priceID string) (string, error) {
	summary, err := s.Summary(ctx, workspaceID, userID)
	if err != nil {
		return "", err
	}
	if summary.Membership != "owner" && summary.Membership != "admin" {
		return "", ErrBillingForbidden
	}
	return s.provider.Checkout(ctx, workspaceID, summary.Subscription.ProviderCustomerID, priceID)
}

func (s *BillingService) Portal(ctx context.Context, workspaceID, userID, returnURL string) (string, error) {
	summary, err := s.Summary(ctx, workspaceID, userID)
	if err != nil {
		return "", err
	}
	if summary.Membership != "owner" || summary.Subscription.ProviderCustomerID == "" {
		return "", ErrBillingForbidden
	}
	return s.provider.Portal(ctx, summary.Subscription.ProviderCustomerID, returnURL)
}

func (s *BillingService) Webhook(ctx context.Context, payload []byte, signature string) (bool, error) {
	event, err := s.provider.ParseWebhook(payload, signature)
	if err != nil {
		return false, fmt.Errorf("verify webhook: %w", err)
	}
	graceEnd := (*time.Time)(nil)
	if event.Status == "past_due" {
		value := time.Now().UTC().Add(72 * time.Hour)
		graceEnd = &value
	}
	return s.repository.ApplyEvent(ctx, event.ID, event.Type, payload, domain.Subscription{
		WorkspaceID: event.WorkspaceID, Provider: "stripe", ProviderCustomerID: event.CustomerID,
		ProviderSubID: event.SubscriptionID, Plan: event.Plan, Status: event.Status,
		CurrentPeriodEnd: event.PeriodEnd, GracePeriodEnd: graceEnd, CancelAtPeriodEnd: event.CancelAtEnd,
	})
}

func entitlements(plan, status string, graceEnd *time.Time) map[string]int64 {
	allowed := status == "trialing" || status == "active" || (status == "past_due" && graceEnd != nil && graceEnd.After(time.Now()))
	if !allowed {
		return map[string]int64{"api_requests": 0, "members": 1}
	}
	switch plan {
	case "pro":
		return map[string]int64{"api_requests": 100000, "members": 10}
	case "team":
		return map[string]int64{"api_requests": 1000000, "members": 100}
	default:
		return map[string]int64{"api_requests": 1000, "members": 2}
	}
}
