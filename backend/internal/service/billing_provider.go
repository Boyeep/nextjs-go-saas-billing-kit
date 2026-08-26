package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Boyeep/nextjs-go-saas-billing-kit/backend/internal/config"
)

type FakeBillingProvider struct{}

func (FakeBillingProvider) Checkout(_ context.Context, workspaceID, _ string, priceID string) (string, error) {
	if workspaceID == "" || priceID == "" {
		return "", errors.New("workspace and price are required")
	}
	return fmt.Sprintf("http://localhost:3000/dashboard/billing?fake_checkout=%s", url.QueryEscape(priceID)), nil
}

func (FakeBillingProvider) Portal(_ context.Context, customerID, returnURL string) (string, error) {
	if customerID == "" {
		return "", errors.New("customer is required")
	}
	return returnURL + "?fake_portal=1", nil
}

func (FakeBillingProvider) ParseWebhook(payload []byte, signature string) (BillingEvent, error) {
	if signature != "fake-valid-signature" {
		return BillingEvent{}, errors.New("invalid fake signature")
	}
	var event BillingEvent
	return event, json.Unmarshal(payload, &event)
}

type StripeBillingProvider struct {
	apiKey        string
	webhookSecret string
	successURL    string
	cancelURL     string
	client        *http.Client
}

func NewBillingProvider(cfg config.BillingConfig) (BillingProvider, error) {
	if cfg.Provider == "fake" {
		return FakeBillingProvider{}, nil
	}
	if cfg.Provider != "stripe" || cfg.StripeAPIKey == "" || cfg.WebhookSecret == "" {
		return nil, errors.New("stripe billing requires STRIPE_SECRET_KEY and STRIPE_WEBHOOK_SECRET")
	}
	return &StripeBillingProvider{apiKey: cfg.StripeAPIKey, webhookSecret: cfg.WebhookSecret, successURL: cfg.SuccessURL, cancelURL: cfg.CancelURL, client: &http.Client{Timeout: 10 * time.Second}}, nil
}

func (p *StripeBillingProvider) Checkout(ctx context.Context, workspaceID, customerID, priceID string) (string, error) {
	values := url.Values{"mode": {"subscription"}, "success_url": {p.successURL}, "cancel_url": {p.cancelURL}, "line_items[0][price]": {priceID}, "line_items[0][quantity]": {"1"}, "subscription_data[metadata][workspace_id]": {workspaceID}}
	if customerID != "" {
		values.Set("customer", customerID)
	}
	return p.createSession(ctx, "/v1/checkout/sessions", values)
}

func (p *StripeBillingProvider) Portal(ctx context.Context, customerID, returnURL string) (string, error) {
	return p.createSession(ctx, "/v1/billing_portal/sessions", url.Values{"customer": {customerID}, "return_url": {returnURL}})
}

func (p *StripeBillingProvider) createSession(ctx context.Context, path string, values url.Values) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.stripe.com"+path, strings.NewReader(values.Encode()))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(p.apiKey, "")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return "", err
	}
	if response.StatusCode >= 300 {
		return "", fmt.Errorf("stripe returned status %d", response.StatusCode)
	}
	var result struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(body, &result); err != nil || result.URL == "" {
		return "", errors.New("stripe response did not include a session URL")
	}
	return result.URL, nil
}

func (p *StripeBillingProvider) ParseWebhook(payload []byte, signature string) (BillingEvent, error) {
	timestamp, signatures, err := parseStripeSignature(signature)
	if err != nil || time.Since(time.Unix(timestamp, 0)).Abs() > 5*time.Minute {
		return BillingEvent{}, errors.New("invalid or expired Stripe signature")
	}
	mac := hmac.New(sha256.New, []byte(p.webhookSecret))
	_, _ = fmt.Fprintf(mac, "%d.%s", timestamp, payload)
	expected := mac.Sum(nil)
	valid := false
	for _, candidate := range signatures {
		decoded, decodeErr := hex.DecodeString(candidate)
		valid = valid || (decodeErr == nil && hmac.Equal(decoded, expected))
	}
	if !valid {
		return BillingEvent{}, errors.New("Stripe signature mismatch")
	}
	var envelope struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Data struct {
			Object struct {
				ID                string            `json:"id"`
				Customer          string            `json:"customer"`
				Status            string            `json:"status"`
				CancelAtPeriodEnd bool              `json:"cancel_at_period_end"`
				CurrentPeriodEnd  int64             `json:"current_period_end"`
				Metadata          map[string]string `json:"metadata"`
			} `json:"object"`
		} `json:"data"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return BillingEvent{}, err
	}
	periodEnd := time.Unix(envelope.Data.Object.CurrentPeriodEnd, 0).UTC()
	return BillingEvent{ID: envelope.ID, Type: envelope.Type, WorkspaceID: envelope.Data.Object.Metadata["workspace_id"], CustomerID: envelope.Data.Object.Customer, SubscriptionID: envelope.Data.Object.ID, Plan: envelope.Data.Object.Metadata["plan"], Status: envelope.Data.Object.Status, PeriodEnd: &periodEnd, CancelAtEnd: envelope.Data.Object.CancelAtPeriodEnd}, nil
}

func parseStripeSignature(header string) (int64, []string, error) {
	var timestamp int64
	var signatures []string
	for _, part := range strings.Split(header, ",") {
		key, value, ok := strings.Cut(strings.TrimSpace(part), "=")
		if !ok {
			continue
		}
		if key == "t" {
			timestamp, _ = strconv.ParseInt(value, 10, 64)
		}
		if key == "v1" {
			signatures = append(signatures, value)
		}
	}
	if timestamp == 0 || len(signatures) == 0 {
		return 0, nil, errors.New("invalid Stripe-Signature header")
	}
	return timestamp, signatures, nil
}
