package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Boyeep/nextjs-go-saas-billing-kit/backend/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrBillingNotFound = errors.New("billing record not found")

type BillingRepository interface {
	ListWorkspaces(context.Context, string) ([]domain.Workspace, error)
	GetSummary(context.Context, string, string) (domain.BillingSummary, error)
	ApplyEvent(context.Context, string, string, []byte, domain.Subscription) (bool, error)
}

func (r *PostgresBillingRepository) ListWorkspaces(ctx context.Context, userID string) ([]domain.Workspace, error) {
	rows, err := r.db.Query(ctx, `SELECT w.id,w.name,w.slug,w.created_at FROM workspaces w JOIN workspace_memberships wm ON wm.workspace_id=w.id WHERE wm.user_id=$1 ORDER BY w.created_at`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Workspace, 0)
	for rows.Next() {
		var item domain.Workspace
		if err := rows.Scan(&item.ID, &item.Name, &item.Slug, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

type PostgresBillingRepository struct{ db *pgxpool.Pool }

func NewBillingRepository(db *pgxpool.Pool) *PostgresBillingRepository {
	return &PostgresBillingRepository{db: db}
}

func (r *PostgresBillingRepository) GetSummary(ctx context.Context, workspaceID, userID string) (domain.BillingSummary, error) {
	var summary domain.BillingSummary
	err := r.db.QueryRow(ctx, `
SELECT w.id, w.name, w.slug, w.created_at, wm.role,
       s.workspace_id, s.provider, COALESCE(s.provider_customer_id, ''),
       COALESCE(s.provider_subscription_id, ''), s.plan, s.status,
       s.current_period_end, s.grace_period_end, s.cancel_at_period_end, s.updated_at
FROM workspaces w
JOIN workspace_memberships wm ON wm.workspace_id = w.id AND wm.user_id = $2
JOIN subscriptions s ON s.workspace_id = w.id
WHERE w.id = $1`, workspaceID, userID).Scan(
		&summary.Workspace.ID, &summary.Workspace.Name, &summary.Workspace.Slug,
		&summary.Workspace.CreatedAt, &summary.Membership,
		&summary.Subscription.WorkspaceID, &summary.Subscription.Provider,
		&summary.Subscription.ProviderCustomerID, &summary.Subscription.ProviderSubID,
		&summary.Subscription.Plan, &summary.Subscription.Status,
		&summary.Subscription.CurrentPeriodEnd, &summary.Subscription.GracePeriodEnd,
		&summary.Subscription.CancelAtPeriodEnd, &summary.Subscription.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return summary, ErrBillingNotFound
	}
	if err != nil {
		return summary, err
	}

	var requests int64
	_ = r.db.QueryRow(ctx, `SELECT COALESCE(SUM(quantity), 0) FROM usage_counters WHERE workspace_id=$1 AND metric='api_requests' AND period_start=date_trunc('month', CURRENT_DATE)::date`, workspaceID).Scan(&requests)
	summary.Usage = map[string]int64{"api_requests": requests}
	return summary, nil
}

func (r *PostgresBillingRepository) ApplyEvent(ctx context.Context, id, eventType string, payload []byte, subscription domain.Subscription) (bool, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = tx.Rollback(ctx) }()
	result, err := tx.Exec(ctx, `INSERT INTO billing_events(provider_event_id,event_type,payload) VALUES($1,$2,$3) ON CONFLICT DO NOTHING`, id, eventType, payload)
	if err != nil || result.RowsAffected() == 0 {
		return false, err
	}
	result, err = tx.Exec(ctx, `UPDATE subscriptions SET provider=$2, provider_customer_id=NULLIF($3,''), provider_subscription_id=NULLIF($4,''), plan=$5, status=$6, current_period_end=$7, grace_period_end=$8, cancel_at_period_end=$9, updated_at=$10 WHERE workspace_id=$1`,
		subscription.WorkspaceID, subscription.Provider, subscription.ProviderCustomerID,
		subscription.ProviderSubID, subscription.Plan, subscription.Status,
		subscription.CurrentPeriodEnd, subscription.GracePeriodEnd,
		subscription.CancelAtPeriodEnd, time.Now().UTC())
	if err != nil {
		return false, err
	}
	if result.RowsAffected() != 1 {
		return false, ErrBillingNotFound
	}
	return true, tx.Commit(ctx)
}
