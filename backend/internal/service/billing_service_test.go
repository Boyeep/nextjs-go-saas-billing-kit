package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Boyeep/nextjs-go-saas-billing-kit/backend/internal/domain"
)

type billingRepoStub struct {
	summary domain.BillingSummary
	events  map[string]bool
	applied domain.Subscription
}

func (r *billingRepoStub) ListWorkspaces(context.Context, string) ([]domain.Workspace, error) {
	return []domain.Workspace{r.summary.Workspace}, nil
}

func (r *billingRepoStub) GetSummary(context.Context, string, string) (domain.BillingSummary, error) {
	return r.summary, nil
}

func (r *billingRepoStub) ApplyEvent(_ context.Context, id, _ string, _ []byte, subscription domain.Subscription) (bool, error) {
	if r.events[id] {
		return false, nil
	}
	r.events[id] = true
	r.applied = subscription
	return true, nil
}

type billingProviderStub struct{ event BillingEvent }

func (billingProviderStub) Checkout(context.Context, string, string, string) (string, error) {
	return "https://checkout.test/session", nil
}
func (billingProviderStub) Portal(context.Context, string, string) (string, error) {
	return "https://portal.test/session", nil
}
func (p billingProviderStub) ParseWebhook([]byte, string) (BillingEvent, error) { return p.event, nil }

func TestSummaryCalculatesActiveEntitlements(t *testing.T) {
	repo := &billingRepoStub{summary: domain.BillingSummary{Membership: "owner", Subscription: domain.Subscription{Plan: "pro", Status: "active"}}}
	summary, err := NewBillingService(repo, billingProviderStub{}).Summary(context.Background(), "workspace", "user")
	if err != nil {
		t.Fatal(err)
	}
	if summary.Entitlements["api_requests"] != 100000 {
		t.Fatalf("unexpected quota: %d", summary.Entitlements["api_requests"])
	}
}

func TestMemberCannotChangeBilling(t *testing.T) {
	repo := &billingRepoStub{summary: domain.BillingSummary{Membership: "member", Subscription: domain.Subscription{Plan: "free", Status: "active"}}}
	_, err := NewBillingService(repo, billingProviderStub{}).Checkout(context.Background(), "workspace", "user", "price")
	if !errors.Is(err, ErrBillingForbidden) {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestWebhookIsIdempotentAndCreatesGracePeriod(t *testing.T) {
	repo := &billingRepoStub{events: map[string]bool{}}
	provider := billingProviderStub{event: BillingEvent{ID: "evt_1", Type: "customer.subscription.updated", WorkspaceID: "ws", Plan: "pro", Status: "past_due"}}
	service := NewBillingService(repo, provider)
	processed, err := service.Webhook(context.Background(), []byte(`{}`), "signature")
	if err != nil || !processed {
		t.Fatalf("first event failed: %v", err)
	}
	if repo.applied.GracePeriodEnd == nil || !repo.applied.GracePeriodEnd.After(time.Now()) {
		t.Fatal("expected grace period")
	}
	processed, err = service.Webhook(context.Background(), []byte(`{}`), "signature")
	if err != nil || processed {
		t.Fatalf("duplicate should be ignored: %v", err)
	}
}

func TestFakeProviderRequiresSignature(t *testing.T) {
	_, err := (FakeBillingProvider{}).ParseWebhook([]byte(`{}`), "wrong")
	if err == nil {
		t.Fatal("expected signature error")
	}
}
