package domain

import "time"

type HealthStatus struct {
	Name      string    `json:"name"`
	Env       string    `json:"env"`
	Version   string    `json:"version"`
	Database  string    `json:"database"`
	Timestamp time.Time `json:"timestamp"`
}

type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	FullName     string     `json:"full_name,omitempty"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	PasswordHash string     `json:"-"`
}

type RegisterInput struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ForgotPasswordInput struct {
	Email string `json:"email"`
}

type ResetPasswordInput struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type VerifyEmailInput struct {
	Token string `json:"token"`
}

type ResendVerificationInput struct {
	Email string `json:"email"`
}

type AuthResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresAt   time.Time `json:"expires_at"`
	User        User      `json:"user"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
}

type Workspace struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}

type Subscription struct {
	WorkspaceID        string     `json:"workspace_id"`
	Provider           string     `json:"provider"`
	ProviderCustomerID string     `json:"provider_customer_id,omitempty"`
	ProviderSubID      string     `json:"provider_subscription_id,omitempty"`
	Plan               string     `json:"plan"`
	Status             string     `json:"status"`
	CurrentPeriodEnd   *time.Time `json:"current_period_end,omitempty"`
	GracePeriodEnd     *time.Time `json:"grace_period_end,omitempty"`
	CancelAtPeriodEnd  bool       `json:"cancel_at_period_end"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

type BillingSummary struct {
	Workspace    Workspace        `json:"workspace"`
	Membership   string           `json:"membership_role"`
	Subscription Subscription     `json:"subscription"`
	Entitlements map[string]int64 `json:"entitlements"`
	Usage        map[string]int64 `json:"usage"`
}
