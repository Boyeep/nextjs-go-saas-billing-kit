package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Boyeep/nextjs-go-saas-billing-kit/backend/internal/service"
)

type BillingHandler struct{ service *service.BillingService }

func NewBillingHandler(billingService *service.BillingService) *BillingHandler {
	return &BillingHandler{service: billingService}
}

func (h *BillingHandler) Workspaces(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		WriteErrorPublic(w, http.StatusUnauthorized, "authentication required")
		return
	}
	items, err := h.service.Workspaces(r.Context(), userID)
	if err != nil {
		WriteErrorPublic(w, http.StatusInternalServerError, "could not load workspaces")
		return
	}
	WriteJSON(w, http.StatusOK, items)
}

func (h *BillingHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		WriteErrorPublic(w, http.StatusUnauthorized, "authentication required")
		return
	}
	summary, err := h.service.Summary(r.Context(), r.PathValue("workspaceID"), userID)
	if errors.Is(err, service.ErrBillingForbidden) {
		WriteErrorPublic(w, http.StatusForbidden, "workspace access denied")
		return
	}
	if err != nil {
		WriteErrorPublic(w, http.StatusInternalServerError, "could not load billing summary")
		return
	}
	WriteJSON(w, http.StatusOK, summary)
}

func (h *BillingHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		WriteErrorPublic(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var input struct {
		PriceID string `json:"price_id"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 64<<10)).Decode(&input); err != nil || input.PriceID == "" {
		WriteErrorPublic(w, http.StatusBadRequest, "price_id is required")
		return
	}
	url, err := h.service.Checkout(r.Context(), r.PathValue("workspaceID"), userID, input.PriceID)
	if errors.Is(err, service.ErrBillingForbidden) {
		WriteErrorPublic(w, http.StatusForbidden, "billing changes require an owner or admin")
		return
	}
	if err != nil {
		WriteErrorPublic(w, http.StatusBadGateway, "could not create checkout session")
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]string{"url": url})
}

func (h *BillingHandler) Portal(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserIDFromContext(r.Context())
	if !ok {
		WriteErrorPublic(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var input struct {
		ReturnURL string `json:"return_url"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 64<<10)).Decode(&input); err != nil || input.ReturnURL == "" {
		WriteErrorPublic(w, http.StatusBadRequest, "return_url is required")
		return
	}
	url, err := h.service.Portal(r.Context(), r.PathValue("workspaceID"), userID, input.ReturnURL)
	if errors.Is(err, service.ErrBillingForbidden) {
		WriteErrorPublic(w, http.StatusForbidden, "customer portal requires the workspace owner")
		return
	}
	if err != nil {
		WriteErrorPublic(w, http.StatusBadGateway, "could not create portal session")
		return
	}
	WriteJSON(w, http.StatusCreated, map[string]string{"url": url})
}

func (h *BillingHandler) Webhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		WriteErrorPublic(w, http.StatusBadRequest, "invalid webhook payload")
		return
	}
	processed, err := h.service.Webhook(r.Context(), payload, r.Header.Get("Stripe-Signature"))
	if err != nil {
		WriteErrorPublic(w, http.StatusBadRequest, "invalid webhook")
		return
	}
	WriteJSON(w, http.StatusOK, map[string]bool{"processed": processed})
}
