package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/auth"
	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/platform/crypto"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"github.com/fakihariefnoto/bloberry/internal/platform/mailer"
)

type usecase struct {
	d auth.Deps
}

func NewUsecase(d auth.Deps) auth.Usecase {
	return &usecase{d: d}
}

var _ auth.Usecase = (*usecase)(nil)

const (
	otpTTL        = 5 * time.Minute
	resetTTL      = 30 * time.Minute
	pairTTL       = 2 * time.Minute
	pendingTTL    = 5 * time.Minute
	maxOtpAttempts = 5
)

func (u *usecase) Signup(ctx context.Context, inviteToken, email, password, displayName, platform string) (*auth.TokenResult, error) {	hash := crypto.HashToken(inviteToken)
	inv, err := u.d.Repo.GetInvitationByTokenHash(ctx, hash)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrInviteInvalid, 400)
	}
	if time.Now().After(inv.ExpiresAt) || inv.AcceptedAt != nil {
		return nil, httpx.NewError(httpx.ErrInviteInvalid, 400)
	}

	usr, err := u.d.Users.GetByEmail(ctx, email)
	if err != nil {
		if !httpx.IsNotFound(err) {
			return nil, err
		}
		ph, err := crypto.HashPassword(password)
		if err != nil {
			return nil, err
		}
		usr = &domain.User{
			Email:        email,
			PasswordHash: &ph,
			DisplayName:  displayName,
		}
		if err := u.d.Users.Insert(ctx, usr); err != nil {
			return nil, err
		}
	}
	if err := u.d.Repo.InsertMembership(ctx, &domain.Membership{
		UserID: usr.ID, TenantID: inv.TenantID, Role: inv.Role,
	}); err != nil {
		return nil, err
	}
	_ = u.d.Repo.MarkInvitationAccepted(ctx, inv.ID)
	return u.issueTokens(ctx, usr, platform)
}

// Activate sets a first-time password for a user who was added as a project
// member without SMTP email (the admin auto-registered them). It only works
// once: if the user already has a password (or was created with one), it
// rejects. Returns the user so the handler can issue tokens.
func (u *usecase) Activate(ctx context.Context, email, password, displayName, platform string) (*auth.TokenResult, error) {
	if email == "" || len(password) < 8 {
		return nil, httpx.NewError(httpx.ErrBadRequest, 400)
	}
	usr, err := u.d.Users.GetByEmail(ctx, email)
	if err != nil {
		if httpx.IsNotFound(err) {
			// No pre-registered account for this email — activation is only
			// for members an admin already added.
			return nil, httpx.NewErrorContent(httpx.ErrActivationInvalid, 403, "No pending account for this email. Ask a project admin to add you first.")
		}
		return nil, err
	}
	if usr.PasswordHash != nil && *usr.PasswordHash != "" {
		return nil, httpx.NewErrorContent(httpx.ErrActivationInvalid, 403, "This account is already activated. Use the forgot-password flow instead.")
	}
	ph, err := crypto.HashPassword(password)
	if err != nil {
		return nil, err
	}
	usr.PasswordHash = &ph
	if displayName != "" {
		usr.DisplayName = displayName
	}
	if err := u.d.Users.Update(ctx, usr); err != nil {
		return nil, err
	}
	return u.issueTokens(ctx, usr, platform)
}

// RegisterMember creates a user with an admin-provided password and adds them
// to a project, bypassing email invitations (used when SMTP is not configured).
func (u *usecase) RegisterMember(ctx context.Context, email, password, displayName, tenantID, role string) (*auth.TokenResult, error) {
	if email == "" || len(password) < 8 {
		return nil, httpx.NewError(httpx.ErrBadRequest, 400)
	}
	usr, err := u.d.Users.GetByEmail(ctx, email)
	if err != nil {
		if !httpx.IsNotFound(err) {
			return nil, err
		}
		ph, err := crypto.HashPassword(password)
		if err != nil {
			return nil, err
		}
		usr = &domain.User{Email: email, PasswordHash: &ph, DisplayName: displayName}
		if err := u.d.Users.Insert(ctx, usr); err != nil {
			return nil, err
		}
	} else if usr.PasswordHash == nil || *usr.PasswordHash == "" {
		ph, err := crypto.HashPassword(password)
		if err != nil {
			return nil, err
		}
		usr.PasswordHash = &ph
		_ = u.d.Users.Update(ctx, usr)
	}
	if err := u.d.Repo.InsertMembership(ctx, &domain.Membership{
		UserID: usr.ID, TenantID: tenantID, Role: role,
	}); err != nil {
		return nil, httpx.NewError(httpx.ErrMemberExists, 409)
	}
	return &auth.TokenResult{User: usr}, nil
}

func (u *usecase) Login(ctx context.Context, email, password, platform string) (*auth.TokenResult, error) {
	usr, err := u.d.Users.GetByEmail(ctx, email)
	authErr := httpx.NewError(httpx.ErrInvalidCredentials, 401)
	if err != nil {
		// constant-time-ish: still do a dummy compare
		if !httpx.IsNotFound(err) {
			return nil, err
		}
		_, _ = crypto.HashPassword(password)
		return nil, authErr
	}
	if usr.PasswordHash == nil {
		return nil, authErr
	}
	ok, err := crypto.VerifyPassword(*usr.PasswordHash, password)
	if err != nil || !ok {
		return nil, authErr
	}
	return u.finishPrimaryLogin(ctx, usr, platform)
}

func (u *usecase) finishPrimaryLogin(ctx context.Context, usr *domain.User, platform string) (*auth.TokenResult, error) {
	if usr.TOTP != nil && usr.TOTP.Enabled {
		pending := crypto.NewToken(24)
		if err := u.d.Redis.Set(ctx, "totp:pending:"+pending, usr.ID, pendingTTL); err != nil {
			return nil, err
		}
		_ = u.d.Users.UpdateLastLogin(ctx, usr.ID)
		return &auth.TokenResult{TotpRequired: true, Pending: pending, User: usr}, nil
	}
	return u.issueTokens(ctx, usr, platform)
}

func (u *usecase) VerifyTotpLogin(ctx context.Context, pending, code string) (*auth.TokenResult, error) {
	userID, err := u.d.Redis.Get(ctx, "totp:pending:"+pending)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrTotpInvalid, 401)
	}
	usr, err := u.d.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrTotpInvalid, 401)
	}
	if usr.TOTP == nil || !usr.TOTP.Enabled {
		return nil, httpx.NewError(httpx.ErrTotpInvalid, 401)
	}
	secret, err := u.d.Envelope.DecryptString(usr.TOTP.SecretEncrypted)
	if err != nil {
		return nil, err
	}
	if crypto.VerifyTOTP(secret, code) {
		_ = u.d.Redis.Del(ctx, "totp:pending:"+pending)
		return u.issueTokens(ctx, usr, "mobile")
	}
	// backup code check
	for i := range usr.TOTP.BackupCodes {
		bc := &usr.TOTP.BackupCodes[i]
		if !bc.Used {
			ok, err := crypto.VerifyPassword(bc.Hash, code)
			if err == nil && ok {
				bc.Used = true
				_ = u.d.Users.Update(ctx, usr)
				_ = u.d.Redis.Del(ctx, "totp:pending:"+pending)
				return u.issueTokens(ctx, usr, "mobile")
			}
		}
	}
	return nil, httpx.NewError(httpx.ErrTotpInvalid, 401)
}

func (u *usecase) Refresh(ctx context.Context, refreshToken string) (*auth.TokenResult, error) {
	userID, platform, err := u.d.Sessions.Get(ctx, refreshToken)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrRefreshInvalid, 401)
	}
	usr, err := u.d.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrRefreshInvalid, 401)
	}
	if err := u.d.Sessions.Revoke(ctx, refreshToken); err != nil {
		return nil, err
	}
	return u.issueTokens(ctx, usr, platform)
}

func (u *usecase) Logout(ctx context.Context, refreshToken string) error {
	return u.d.Sessions.Revoke(ctx, refreshToken)
}

func (u *usecase) ForgotPassword(ctx context.Context, email string) error {
	usr, err := u.d.Users.GetByEmail(ctx, email)
	if err != nil {
		return nil // identical response whether or not email exists
	}
	token := crypto.NewToken(24)
	key := "reset:" + crypto.HashToken(token)
	if err := u.d.Redis.Set(ctx, key, usr.ID, resetTTL); err != nil {
		return err
	}
	url := u.d.BaseURL + "/reset-password?token=" + token
	_ = u.d.Mailer.Send(ctx, email, "Reset your Bloberry password",
		"Reset link (valid 30 minutes). If you didn't ask for this, ignore it.\n\n"+url,
		mailer.Render("reset", map[string]string{"email": email, "url": url}))
	return nil
}

func (u *usecase) ResetPassword(ctx context.Context, token, newPassword string) error {
	key := "reset:" + crypto.HashToken(token)
	userID, err := u.d.Redis.Get(ctx, key)
	if err != nil {
		return httpx.NewError(httpx.ErrResetTokenInvalid, 400)
	}
	_ = u.d.Redis.Del(ctx, key)
	usr, err := u.d.Users.GetByID(ctx, userID)
	if err != nil {
		return httpx.NewError(httpx.ErrResetTokenInvalid, 400)
	}
	hash, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}
	usr.PasswordHash = &hash
	if err := u.d.Users.Update(ctx, usr); err != nil {
		return err
	}
	// invalidate every existing session
	_ = u.d.Sessions.RevokeAll(ctx, usr.ID)
	return nil
}

func (u *usecase) RequestOTP(ctx context.Context, email string) error {
	cnt, err := u.d.Redis.Incr(ctx, "otp:attempts:"+email)
	if err != nil {
		return err
	}
	_ = u.d.Redis.Expire(ctx, "otp:attempts:"+email, time.Hour)
	if cnt > 5 {
		return httpx.NewError(httpx.ErrOtpRateLimited, 429)
	}
	code := fmt.Sprintf("%06d", crypto.Random6())
	key := "otp:login:" + email
	payload, _ := json.Marshal(map[string]interface{}{"code_hash": crypto.HashToken(code), "attempts": 0})
	if err := u.d.Redis.Set(ctx, key, string(payload), otpTTL); err != nil {
		return err
	}
	_ = u.d.Mailer.Send(ctx, email, "Your Bloberry login code",
		"Your one-time code is: "+code+" (valid 5 minutes)",
		mailer.Render("otp", map[string]string{"code": code}))
	return nil
}

func (u *usecase) VerifyOTP(ctx context.Context, email, code, platform string) (*auth.TokenResult, error) {
	val, err := u.d.Redis.Get(ctx, "otp:login:"+email)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrOtpInvalid, 401)
	}
	var payload struct {
		CodeHash string `json:"code_hash"`
		Attempts int    `json:"attempts"`
	}
	_ = json.Unmarshal([]byte(val), &payload)
	if payload.Attempts >= maxOtpAttempts || crypto.HashToken(code) != payload.CodeHash {
		_ = u.d.Redis.Del(ctx, "otp:login:"+email)
		return nil, httpx.NewError(httpx.ErrOtpInvalid, 401)
	}
	_ = u.d.Redis.Del(ctx, "otp:login:"+email)
	usr, err := u.d.Users.GetByEmail(ctx, email)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrOtpInvalid, 401)
	}
	return u.finishPrimaryLogin(ctx, usr, platform)
}

func (u *usecase) GoogleLogin(ctx context.Context, idToken, platform string) (*auth.TokenResult, error) {
	gid, err := u.d.Google.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrOauthInvalid, 401)
	}
	// find identity
	usr, err := u.d.Users.GetByEmail(ctx, gid.Email)
	if err != nil {
		if !httpx.IsNotFound(err) {
			return nil, err
		}
		return nil, httpx.NewErrorContent(httpx.ErrNoInvitation, 403, "No account for this identity")
	}
	// link identity
	found := false
	for _, id := range usr.OAuthIdentities {
		if id.Provider == "google" && id.ProviderUserID == gid.ProviderUserID {
			found = true
			break
		}
	}
	if !found && gid.EmailVerified {
		usr.OAuthIdentities = append(usr.OAuthIdentities, domain.OAuthID{
			Provider: "google", ProviderUserID: gid.ProviderUserID, EmailAtLink: gid.Email, CreatedAt: time.Now().UTC(),
		})
		_ = u.d.Users.Update(ctx, usr)
	}
	return u.finishPrimaryLogin(ctx, usr, platform)
}

func (u *usecase) ProvisionTOTP(ctx context.Context, userID string) (*auth.TOTPProvision, error) {
	usr, err := u.d.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if usr.TOTP != nil && usr.TOTP.Enabled {
		return nil, httpx.NewError(httpx.ErrConflict, 409)
	}
	secret := crypto.GenerateTOTPSecret()
	enc, err := u.d.Envelope.EncryptString(secret)
	if err != nil {
		return nil, err
	}
	if usr.TOTP == nil {
		usr.TOTP = &domain.TOTPConfig{}
	}
	usr.TOTP.SecretEncrypted = enc
	usr.TOTP.Enabled = false
	if err := u.d.Users.Update(ctx, usr); err != nil {
		return nil, err
	}
	issuer := "bloberry"
	url := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s", issuer, usr.Email, secret, issuer)
	return &auth.TOTPProvision{Secret: secret, OtpAuthURL: url}, nil
}

func (u *usecase) EnableTOTP(ctx context.Context, userID, code string) ([]string, error) {
	usr, err := u.d.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if usr.TOTP == nil || usr.TOTP.SecretEncrypted == "" {
		return nil, httpx.NewError(httpx.ErrTotpInvalid, 400)
	}
	secret, err := u.d.Envelope.DecryptString(usr.TOTP.SecretEncrypted)
	if err != nil {
		return nil, err
	}
	if !crypto.VerifyTOTP(secret, code) {
		return nil, httpx.NewError(httpx.ErrTotpInvalid, 400)
	}
	now := time.Now().UTC()
	usr.TOTP.Enabled = true
	usr.TOTP.EnabledAt = &now
	backup := crypto.GenerateBackupCodes(10)
	codes := make([]domain.BackupCode, 0, len(backup))
	for _, c := range backup {
		h, _ := crypto.HashPassword(c)
		codes = append(codes, domain.BackupCode{Hash: h})
	}
	usr.TOTP.BackupCodes = codes
	if err := u.d.Users.Update(ctx, usr); err != nil {
		return nil, err
	}
	return backup, nil
}

func (u *usecase) IssuePairToken(ctx context.Context, userID string) (string, error) {
	token := crypto.NewToken(24)
	key := "pair:" + token
	if err := u.d.Redis.Set(ctx, key, userID, pairTTL); err != nil {
		return "", err
	}
	return "bloberry://pair/" + token, nil
}

func (u *usecase) VerifyPairToken(ctx context.Context, payload, platform string) (*auth.TokenResult, error) {
	const prefix = "bloberry://pair/"
	if len(payload) <= len(prefix) {
		return nil, httpx.NewErrorContent(httpx.ErrPairInvalid, 400, "This code is no longer valid — refresh it")
	}
	token := payload[len(prefix):]
	userID, err := u.d.Redis.Get(ctx, "pair:"+token)
	if err != nil {
		return nil, httpx.NewErrorContent(httpx.ErrPairInvalid, 400, "This code is no longer valid — refresh it")
	}
	_ = u.d.Redis.Del(ctx, "pair:"+token) // single-use
	usr, err := u.d.Users.GetByID(ctx, userID)
	if err != nil {
		return nil, httpx.NewError(httpx.ErrPairInvalid, 400)
	}
	return u.issueTokens(ctx, usr, platform)
}

func (u *usecase) IssueConfigPayload(ctx context.Context, userID string) (string, error) {
	// Signed import-window payload (domains.md §4.9). The server signs with
	// its HMAC secret; encryption happens client-side with the passphrase.
	payload := map[string]interface{}{
		"user_id":                userID,
		"import_window_expires":  time.Now().Add(24 * time.Hour).Unix(),
		"server":                 "bloberry",
	}
	raw, _ := json.Marshal(payload)
	sig := crypto.HashToken("config:" + string(raw))
	sealed := map[string]interface{}{
		"payload": string(raw),
		"signature": sig,
	}
	out, _ := json.Marshal(sealed)
	return string(out), nil
}

// --- helpers ---

func (u *usecase) issueTokens(ctx context.Context, usr *domain.User, platform string) (*auth.TokenResult, error) {
	accessTTL := int64(720 * 3600)   // mobile 720h
	refreshTTL := int64(2160 * 3600) // mobile 2160h
	if platform == "web" {
		accessTTL = int64(48 * 3600)
		refreshTTL = int64(144 * 3600)
	}
	access, err := u.d.Tokens.Issue(usr.ID, time.Duration(accessTTL)*time.Second)
	if err != nil {
		return nil, err
	}
	refresh, err := u.d.Sessions.Create(ctx, usr.ID, platform, refreshTTL)
	if err != nil {
		return nil, err
	}
	_ = u.d.Users.UpdateLastLogin(ctx, usr.ID)
	return &auth.TokenResult{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    accessTTL,
		User:         usr,
	}, nil
}

var _ = errors.New
