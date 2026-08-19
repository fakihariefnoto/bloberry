package auth

import (
	"context"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/user"
)

// Repository is the auth domain's persistence boundary. It reads users via
// user.Reader and owns invitation + membership writes.
type Repository interface {
	GetInvitationByTokenHash(ctx context.Context, hash string) (*domain.Invitation, error)
	MarkInvitationAccepted(ctx context.Context, id string) error
	InsertMembership(ctx context.Context, m *domain.Membership) error
	GetMembership(ctx context.Context, userID, tenantID string) (*domain.Membership, error)
	InsertTenant(ctx context.Context, t *domain.Tenant) error
	GetTenantBySlug(ctx context.Context, slug string) (*domain.Tenant, error)
	InsertRootFolder(ctx context.Context, f *domain.Folder) error
}

// TokenResult is what the handler returns to clients.
type TokenResult struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int64         `json:"expires_in"`
	User         *domain.User  `json:"user,omitempty"`
	TotpRequired bool          `json:"totp_required,omitempty"`
	Pending      string        `json:"pending,omitempty"`
}

// Usecase is the auth service.
type Usecase interface {
	Signup(ctx context.Context, inviteToken, email, password, displayName, platform string) (*TokenResult, error)
	Login(ctx context.Context, email, password, platform string) (*TokenResult, error)
	Activate(ctx context.Context, email, password, displayName, platform string) (*TokenResult, error)
	RegisterMember(ctx context.Context, email, password, displayName, tenantID, role string) (*TokenResult, error)
	VerifyTotpLogin(ctx context.Context, pending, code string) (*TokenResult, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenResult, error)
	Logout(ctx context.Context, refreshToken string) error
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	RequestOTP(ctx context.Context, email string) error
	VerifyOTP(ctx context.Context, email, code, platform string) (*TokenResult, error)
	GoogleLogin(ctx context.Context, idToken, platform string) (*TokenResult, error)
	ProvisionTOTP(ctx context.Context, userID string) (*TOTPProvision, error)
	EnableTOTP(ctx context.Context, userID, code string) ([]string, error)
	IssuePairToken(ctx context.Context, userID string) (string, error)
	VerifyPairToken(ctx context.Context, payload, platform string) (*TokenResult, error)
	IssueConfigPayload(ctx context.Context, userID string) (string, error)
}

type TOTPProvision struct {
	Secret     string `json:"secret"`
	OtpAuthURL string `json:"otpauth_url"`
}

// Dependencies the usecase needs; satisfied by the wiring in main.
type Deps struct {
	Repo     Repository
	Users    user.Repository
	Sessions SessionStore
	Tokens   TokenIssuer
	Mailer   Mailer
	Google   GoogleVerifier
	Envelope CredentialEnvelope
	Redis    RedisStore
}

type RedisStore interface {
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Incr(ctx context.Context, key string) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
}

type SessionStore interface {
	Create(ctx context.Context, userID, platform string, ttlSeconds int64) (string, error)
	Rotate(ctx context.Context, sid, userID, platform string, ttlSeconds int64) (string, error)
	Get(ctx context.Context, sid string) (userID string, platform string, err error)
	Revoke(ctx context.Context, sid string) error
	RevokeAll(ctx context.Context, userID string) error
}

type TokenIssuer interface {
	Issue(userID string, ttl time.Duration) (string, error)
}

type Mailer interface {
	Send(ctx context.Context, to, subject, text string) error
}

type GoogleVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*GoogleIdentity, error)
}

type GoogleIdentity struct {
	ProviderUserID string
	Email          string
	EmailVerified  bool
}

// CredentialEnvelope encrypts the TOTP secret at rest (domains.md §4.10).
type CredentialEnvelope interface {
	EncryptString(plaintext string) (string, error)
	DecryptString(encoded string) (string, error)
}
