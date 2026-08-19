package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fakihariefnoto/bloberry/internal/admin"
	"github.com/fakihariefnoto/bloberry/internal/apikey"
	"github.com/fakihariefnoto/bloberry/internal/audit"
	"github.com/fakihariefnoto/bloberry/internal/auth"
	"github.com/fakihariefnoto/bloberry/internal/authz"
	"github.com/fakihariefnoto/bloberry/internal/domain"
	"github.com/fakihariefnoto/bloberry/internal/folder"
	"github.com/fakihariefnoto/bloberry/internal/grant"
	"github.com/fakihariefnoto/bloberry/internal/job"
	"github.com/fakihariefnoto/bloberry/internal/object"
	server "github.com/fakihariefnoto/bloberry/internal/platform/api"
	"github.com/fakihariefnoto/bloberry/internal/platform/httpx"
	"github.com/fakihariefnoto/bloberry/internal/share"
	"github.com/fakihariefnoto/bloberry/internal/setup"
	"github.com/fakihariefnoto/bloberry/internal/storage"
	"github.com/fakihariefnoto/bloberry/internal/tenant"
	"github.com/fakihariefnoto/bloberry/internal/usage"
	"github.com/fakihariefnoto/bloberry/internal/user"
)

// Handler implements server.ServerInterface by delegating to the domain usecases.
type Handler struct {
	Auth    auth.Usecase
	Users   user.Usecase
	Tenants tenant.Usecase
	Folders folder.Usecase
	Objects object.Usecase
	Shares  share.Usecase
	APIKeys apikey.Usecase
	Grants  grant.Usecase
	Jobs    job.Usecase
	Usage   usage.Usecase
	Audit   audit.Usecase
	Admin   admin.Usecase
	Setup   setup.Usecase
	// Storage registry for the disk driver's raw HMAC endpoint + health.
	Storage  object.Registry
	Envelope interface {
		Decrypt(ct []byte) ([]byte, error)
	}
}

var _ server.ServerInterface = (*Handler)(nil)

func principalOf(r *http.Request) *authz.Principal { return httpx.PrincipalFrom(r.Context()) }
func tenantOf(r *http.Request) string {
	if p := principalOf(r); p != nil {
		return p.TenantID
	}
	return ""
}

func decodeBody(r *http.Request, v interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// --- helpers to respond with the envelope ---

func data(w http.ResponseWriter, status int, d interface{}) {
	httpx.WriteJSON(w, status, httpx.Envelope{Data: d})
}

func requireAdmin(w http.ResponseWriter, r *http.Request) bool {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	if !p.IsPlatformAdmin {
		httpx.Error(w, http.StatusForbidden, "forbidden")
		return false
	}
	return true
}

// requireTenantAdminOrPlatform allows the platform admin, or a user whose
// principal is an admin/owner of the given tenant. Used for member management
// and tenant settings: a platform admin can manage any project's members, a
// project admin/owner only their own.
func requireTenantAdminOrPlatform(w http.ResponseWriter, r *http.Request, tenantID string) bool {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return false
	}
	if p.IsPlatformAdmin {
		return true
	}
	if p.TenantID == tenantID && (p.Role == authz.RoleTenantAdmin || p.Role == authz.RoleTenantOwner) {
		return true
	}
	httpx.Error(w, http.StatusForbidden, "forbidden")
	return false
}

// --- health ---

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	data(w, http.StatusOK, map[string]interface{}{"status": "ok"})
}

func (h *Handler) SetupStatus(w http.ResponseWriter, r *http.Request) {
	st, err := h.Setup.Status(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, st)
}

func (h *Handler) RunSetup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
		TenantName  string `json:"tenant_name"`
		TenantSlug  string `json:"tenant_slug"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	if err := h.Setup.Run(r.Context(), req.Email, req.Password, req.DisplayName, req.TenantName, req.TenantSlug); err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, map[string]interface{}{"setup_complete": true})
}

// --- auth ---

type signupReq struct {
	InviteToken string `json:"invite_token"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	Platform    string `json:"platform"`
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) {
	var req signupReq
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	res, err := h.Auth.Signup(r.Context(), req.InviteToken, req.Email, req.Password, req.DisplayName, orPlatform(req.Platform))
	writeTokens(w, res, err)
}

// Activate sets the first password for a member who was added without SMTP
// email (auto-registered pending activation). First-time only — the same email
// cannot activate twice.
func (h *Handler) ActivateAccount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
		Platform    string `json:"platform"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	res, err := h.Auth.Activate(r.Context(), req.Email, req.Password, req.DisplayName, orPlatform(req.Platform))
	writeTokens(w, res, err)
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Platform string `json:"platform"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	res, err := h.Auth.Login(r.Context(), req.Email, req.Password, orPlatform(req.Platform))
	writeTokens(w, res, err)
}

type totpLoginReq struct {
	Pending string `json:"pending"`
	Code    string `json:"code"`
}

func (h *Handler) VerifyTotpLogin(w http.ResponseWriter, r *http.Request) {
	var req totpLoginReq
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	res, err := h.Auth.VerifyTotpLogin(r.Context(), req.Pending, req.Code)
	writeTokens(w, res, err)
}

type refreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshReq
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	res, err := h.Auth.Refresh(r.Context(), req.RefreshToken)
	writeTokens(w, res, err)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req refreshReq
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	if err := h.Auth.Logout(r.Context(), req.RefreshToken); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	if err := h.Auth.ForgotPassword(r.Context(), req.Email); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.MessageEnvelope(w, http.StatusOK, httpx.Message{Code: "password_reset_sent", Content: "If that email exists, a reset link is on its way"})
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	if err := h.Auth.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, nil)
}

func (h *Handler) RequestOtp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	if err := h.Auth.RequestOTP(r.Context(), req.Email); err != nil {
		httpx.WriteError(w, err)
		return
	}
	httpx.MessageEnvelope(w, http.StatusOK, httpx.Message{Code: "otp_sent"})
}

func (h *Handler) VerifyOtp(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Code     string `json:"code"`
		Platform string `json:"platform"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	res, err := h.Auth.VerifyOTP(r.Context(), req.Email, req.Code, orPlatform(req.Platform))
	writeTokens(w, res, err)
}

func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDToken  string `json:"id_token"`
		Platform string `json:"platform"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	res, err := h.Auth.GoogleLogin(r.Context(), req.IDToken, orPlatform(req.Platform))
	writeTokens(w, res, err)
}

func (h *Handler) ProvisionTotp(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	prov, err := h.Auth.ProvisionTOTP(r.Context(), p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, prov)
}

func (h *Handler) EnableTotp(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		Code string `json:"code"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	codes, err := h.Auth.EnableTOTP(r.Context(), p.ID, req.Code)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, map[string]interface{}{"backup_codes": codes})
}

func (h *Handler) IssuePairToken(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	payload, err := h.Auth.IssuePairToken(r.Context(), p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, map[string]interface{}{"qr_payload": payload})
}

func (h *Handler) VerifyPairToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Payload  string `json:"payload"`
		Platform string `json:"platform"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	res, err := h.Auth.VerifyPairToken(r.Context(), req.Payload, orPlatform(req.Platform))
	writeTokens(w, res, err)
}

func (h *Handler) IssueConfigFile(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	payload, err := h.Auth.IssueConfigPayload(r.Context(), p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, map[string]interface{}{"payload": payload})
}

func writeTokens(w http.ResponseWriter, res *auth.TokenResult, err error) {
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, res)
}

func orPlatform(p string) string {
	if p == "" {
		return "mobile"
	}
	return p
}

// --- users ---

func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	u, err := h.Users.GetProfile(r.Context(), p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, u)
}

func (h *Handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		DisplayName          *string `json:"display_name"`
		Locale               *string `json:"locale"`
		NotificationsEnabled *bool   `json:"notifications_enabled"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	u, err := h.Users.UpdateProfile(r.Context(), p.ID, req.DisplayName, req.Locale, req.NotificationsEnabled)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, u)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	if err := h.Users.ChangePassword(r.Context(), p.ID, req.CurrentPassword, req.NewPassword); err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, nil)
}

// --- tenants ---

func (h *Handler) ListTenants(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	ts, err := h.Tenants.ListForUser(r.Context(), p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, ts)
}

func (h *Handler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if !p.IsPlatformAdmin {
		httpx.Error(w, http.StatusForbidden, "forbidden")
		return
	}
	var req struct {
		Name             string `json:"name"`
		Slug             string `json:"slug"`
		QuotaBytes       int64  `json:"quota_bytes"`
		QuotaObjects     int64  `json:"quota_objects"`
		DefaultBackendID string `json:"default_storage_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	t, err := h.Tenants.Create(r.Context(), req.Name, req.Slug, req.QuotaBytes, req.QuotaObjects, req.DefaultBackendID, p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, t)
}

func (h *Handler) GetTenant(w http.ResponseWriter, r *http.Request, tenantId server.TenantId) {
	if !requireTenantAdminOrPlatform(w, r, string(tenantId)) {
		return
	}
	t, err := h.Tenants.Get(r.Context(), string(tenantId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, t)
}

func (h *Handler) UpdateTenant(w http.ResponseWriter, r *http.Request, tenantId server.TenantId) {
	if !requireTenantAdminOrPlatform(w, r, string(tenantId)) {
		return
	}
	var req struct {
		Name             *string  `json:"name"`
		QuotaBytes       *int64   `json:"quota_bytes"`
		QuotaObjects     *int64   `json:"quota_objects"`
		DefaultBackendID *string  `json:"default_storage_id"`
		StorageEngines   *[]string `json:"storage_engines"`
		Status           *string  `json:"status"`
		UploadPolicy     *domain.UploadPolicy `json:"upload_policy"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	t, err := h.Tenants.Update(r.Context(), string(tenantId), tenant.Update{
		Name: req.Name, QuotaBytes: req.QuotaBytes, QuotaObjects: req.QuotaObjects,
		DefaultBackendID: req.DefaultBackendID, StorageEngines: req.StorageEngines, Status: req.Status,
		UploadPolicy: req.UploadPolicy,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, t)
}

func (h *Handler) ListMembers(w http.ResponseWriter, r *http.Request, tenantId server.TenantId) {
	if !requireTenantAdminOrPlatform(w, r, string(tenantId)) {
		return
	}
	ms, err := h.Tenants.ListMembers(r.Context(), string(tenantId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	// Enrich each membership with the user's email + display name so the UI
	// can render "Jane (jane@acme.com)" instead of an opaque user id.
	type memberView struct {
		domain.Membership
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	}
	out := make([]memberView, 0, len(ms))
	for _, m := range ms {
		mv := memberView{Membership: m}
		if u, err := h.Users.GetProfile(r.Context(), m.UserID); err == nil {
			mv.Email = u.Email
			mv.DisplayName = u.DisplayName
		}
		out = append(out, mv)
	}
	data(w, http.StatusOK, out)
}

func (h *Handler) AddMember(w http.ResponseWriter, r *http.Request, tenantId server.TenantId) {
	if !requireTenantAdminOrPlatform(w, r, string(tenantId)) {
		return
	}
	var req struct {
		UserID   string `json:"user_id"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Password string `json:"password"`
		// Method: "invite" (default, requires SMTP) | "password" (admin sets a
		// password, user is auto-registered) | "activation" (no email server;
		// user activates on first login by email + password).
		Method string `json:"method"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	if req.Email == "" && req.UserID == "" {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}

	// Method "password": auto-register (create user + password + membership)
	// without any email. This is the no-SMTP path.
	if req.Method == "password" {
		if req.Password == "" {
			httpx.Error(w, http.StatusBadRequest, "password_required")
			return
		}
		res, err := h.Auth.RegisterMember(r.Context(), req.Email, req.Password, "", string(tenantId), req.Role)
		if err != nil {
			httpx.WriteError(w, err)
			return
		}
		data(w, http.StatusCreated, map[string]interface{}{
			"method": "password", "user_id": res.User.ID, "email": res.User.Email,
		})
		return
	}

	// Method "activation": no SMTP configured — register a password-less user
	// that must activate once (email + password) before first login.
	if req.Method == "activation" {
		ph := ""
		usr, err := h.Users.GetByEmail(r.Context(), req.Email)
		if err == nil {
			if usr.PasswordHash != nil {
				httpx.Error(w, http.StatusConflict, "already_activated")
				return
			}
		} else if httpx.IsNotFound(err) {
			usr, err = h.Users.CreatePending(r.Context(), req.Email)
			if err != nil {
				httpx.WriteError(w, err)
				return
			}
		} else {
			httpx.WriteError(w, err)
			return
		}
		_ = ph
		if err := h.Tenants.AddMember(r.Context(), string(tenantId), usr.ID, req.Role); err != nil {
			httpx.WriteError(w, err)
			return
		}
		data(w, http.StatusCreated, map[string]interface{}{"method": "activation", "user_id": usr.ID, "email": usr.Email})
		return
	}

	// Default "invite" method: user must already exist, or SMTP must be
	// configured to deliver the invitation email.
	userID := req.UserID
	if userID == "" && req.Email != "" {
		u, err := h.Users.GetByEmail(r.Context(), req.Email)
		if err != nil {
			httpx.Error(w, http.StatusNotFound, "user_not_found")
			return
		}
		userID = u.ID
	}
	if userID == "" {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	if err := h.Tenants.AddMember(r.Context(), string(tenantId), userID, req.Role); err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, map[string]interface{}{"method": "invite", "user_id": userID})
}

func (h *Handler) RemoveMember(w http.ResponseWriter, r *http.Request, tenantId server.TenantId, membershipId string) {
	if !requireTenantAdminOrPlatform(w, r, string(tenantId)) {
		return
	}
	if err := h.Tenants.RemoveMember(r.Context(), string(tenantId), membershipId); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateMember(w http.ResponseWriter, r *http.Request, tenantId server.TenantId, membershipId string) {
	if !requireTenantAdminOrPlatform(w, r, string(tenantId)) {
		return
	}
	var req struct {
		Role string `json:"role"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	if err := h.Tenants.UpdateMemberRole(r.Context(), string(tenantId), membershipId, req.Role); err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, nil)
}

func (h *Handler) ListInvitations(w http.ResponseWriter, r *http.Request, tenantId server.TenantId) {
	if !requireTenantAdminOrPlatform(w, r, string(tenantId)) {
		return
	}
	is, err := h.Tenants.ListInvitations(r.Context(), string(tenantId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, is)
}

func (h *Handler) CreateInvitation(w http.ResponseWriter, r *http.Request, tenantId server.TenantId) {
	if !requireTenantAdminOrPlatform(w, r, string(tenantId)) {
		return
	}
	p := principalOf(r)
	var req struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	inv, err := h.Tenants.CreateInvitation(r.Context(), string(tenantId), req.Email, req.Role, p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, inv)
}

// --- folders ---

func (h *Handler) CreateFolder(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		Name     string `json:"name"`
		ParentID string `json:"parent_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	parentID := req.ParentID
	if parentID == "" || parentID == "root" {
		root, err := h.Folders.GetRoot(r.Context(), p.TenantID)
		if err == nil {
			parentID = root.ID
		}
	}
	f, err := h.Folders.Create(r.Context(), p.TenantID, parentID, req.Name)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, f)
}

func (h *Handler) GetFolder(w http.ResponseWriter, r *http.Request, folderId server.FolderId) {
	p := principalOf(r)
	f, err := h.Folders.Get(r.Context(), p.TenantID, string(folderId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, f)
}

func (h *Handler) RenameFolder(w http.ResponseWriter, r *http.Request, folderId server.FolderId) {
	p := principalOf(r)
	var req struct {
		Name         string                `json:"name"`
		UploadPolicy *domain.UploadPolicy  `json:"upload_policy"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	var f *domain.Folder
	if req.UploadPolicy != nil {
		pol, err := h.Folders.SetPolicy(r.Context(), p.TenantID, string(folderId), req.UploadPolicy)
		if err != nil {
			httpx.WriteError(w, err)
			return
		}
		f = pol
	}
	if req.Name != "" {
		renamed, err := h.Folders.Rename(r.Context(), p.TenantID, string(folderId), req.Name)
		if err != nil {
			httpx.WriteError(w, err)
			return
		}
		f = renamed
	}
	if f == nil {
		f, _ = h.Folders.Get(r.Context(), p.TenantID, string(folderId))
	}
	data(w, http.StatusOK, f)
}

func (h *Handler) DeleteFolder(w http.ResponseWriter, r *http.Request, folderId server.FolderId) {
	p := principalOf(r)
	if err := h.Folders.Delete(r.Context(), p.TenantID, string(folderId)); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) MoveFolder(w http.ResponseWriter, r *http.Request, folderId server.FolderId) {
	p := principalOf(r)
	var req struct {
		TargetParentID string `json:"target_parent_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	f, err := h.Folders.Move(r.Context(), p.TenantID, string(folderId), req.TargetParentID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, f)
}

func (h *Handler) ListFolderChildren(w http.ResponseWriter, r *http.Request, folderId server.FolderId) {
	p := principalOf(r)
	var fid string
	if string(folderId) == "root" {
		root, err := h.Folders.GetRoot(r.Context(), p.TenantID)
		if err != nil {
			httpx.WriteError(w, err)
			return
		}
		fid = root.ID
	} else {
		fid = string(folderId)
	}
	folders, err := h.Folders.ListChildren(r.Context(), p.TenantID, strPtr(fid))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	storageID := r.URL.Query().Get("storage_id")
	objects, err := h.Objects.ListByFolder(r.Context(), p.TenantID, fid, storageID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	enriched := enrichObjects(r.Context(), h, objects)
	data(w, http.StatusOK, map[string]interface{}{"folders": folders, "objects": enriched})
}

// enrichObjects appends the storage backend's name + driver to each object so
// the UI can show which account/vendor holds the file (PRD: storage-agnostic,
// but the dashboard must still show where bytes actually live).
func enrichObjects(ctx context.Context, h *Handler, objects []domain.Object) []map[string]interface{} {
	cache := map[string]*domain.StorageBackend{}
	out := make([]map[string]interface{}, 0, len(objects))
	for i := range objects {
		o := &objects[i]
		m := map[string]interface{}{}
		b, _ := json.Marshal(o)
		_ = json.Unmarshal(b, &m)
		be, ok := cache[o.BackendID]
		if !ok {
			b, err := h.Admin.GetBackend(ctx, o.BackendID)
			if err == nil {
				be = b
				cache[o.BackendID] = b
			}
		}
		if be != nil {
			m["backend_name"] = be.Name
			m["backend_driver"] = be.Driver
		}
		out = append(out, m)
	}
	return out
}

func (h *Handler) GetFolderTree(w http.ResponseWriter, r *http.Request, folderId server.FolderId) {
	p := principalOf(r)
	_, err := h.Folders.Get(r.Context(), p.TenantID, string(folderId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	// v1: flat list under this folder.
	folders, err := h.Folders.ListChildren(r.Context(), p.TenantID, strPtr(string(folderId)))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, map[string]interface{}{"folders": folders})
}

// --- objects ---

type presignReq struct {
	FolderID    string `json:"folder_id"`
	Name        string `json:"name"`
	BackendID   string `json:"storage_id"`
	Overwrite   bool   `json:"overwrite"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
}

func (h *Handler) PresignPut(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req presignReq
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	res, err := h.Objects.PresignPut(r.Context(), p.TenantID, req.FolderID, req.Name, req.BackendID, req.Overwrite, req.Size, req.ContentType, p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, res)
}

func (h *Handler) CompleteUpload(w http.ResponseWriter, r *http.Request, fileId server.FileId) {
	p := principalOf(r)
	var req struct {
		ETag string `json:"etag"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	obj, err := h.Objects.Complete(r.Context(), p.TenantID, string(fileId), req.ETag)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, obj)
}

func (h *Handler) DirectUpload(w http.ResponseWriter, r *http.Request, params server.DirectUploadParams) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if params.Name == "" || params.FolderId == "" {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	contentType := ""
	if params.ContentType != nil {
		contentType = *params.ContentType
	}
	backendID := r.URL.Query().Get("backend_id")
	obj, err := h.Objects.DirectUpload(r.Context(), p.TenantID, params.FolderId, params.Name, backendID, contentType, r.Body, r.ContentLength, p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, obj)
}

func (h *Handler) MultipartInit(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	var req presignReq
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	res, err := h.Objects.MultipartInit(r.Context(), p.TenantID, req.FolderID, req.Name, req.BackendID, req.Overwrite, req.Size, req.ContentType, p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, res)
}

func (h *Handler) MultipartPresignPart(w http.ResponseWriter, r *http.Request, fileId server.FileId) {
	p := principalOf(r)
	var req struct {
		Part int `json:"part"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	res, err := h.Objects.MultipartPresignPart(r.Context(), p.TenantID, string(fileId), req.Part)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, res)
}

func (h *Handler) MultipartComplete(w http.ResponseWriter, r *http.Request, fileId server.FileId) {
	p := principalOf(r)
	var req struct {
		Parts []storage.Part `json:"parts"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	obj, err := h.Objects.MultipartComplete(r.Context(), p.TenantID, string(fileId), req.Parts)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, obj)
}

func (h *Handler) MultipartStatus(w http.ResponseWriter, r *http.Request, fileId server.FileId) {
	p := principalOf(r)
	mp, err := h.Objects.MultipartStatus(r.Context(), p.TenantID, string(fileId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, mp)
}

func (h *Handler) StatObject(w http.ResponseWriter, r *http.Request, fileId server.FileId) {
	p := principalOf(r)
	obj, err := h.Objects.Get(r.Context(), p.TenantID, string(fileId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, obj)
}

func (h *Handler) DeleteObject(w http.ResponseWriter, r *http.Request, fileId server.FileId) {
	p := principalOf(r)
	if err := h.Objects.Delete(r.Context(), p.TenantID, string(fileId)); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) MoveObject(w http.ResponseWriter, r *http.Request, fileId server.FileId) {
	p := principalOf(r)
	var req struct {
		TargetFolderID string `json:"target_folder_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	obj, err := h.Objects.Move(r.Context(), p.TenantID, string(fileId), req.TargetFolderID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, obj)
}

func (h *Handler) SetVisibility(w http.ResponseWriter, r *http.Request, fileId server.FileId) {
	p := principalOf(r)
	var req struct {
		Visibility string `json:"visibility"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	obj, err := h.Objects.SetVisibility(r.Context(), p.TenantID, string(fileId), req.Visibility)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, obj)
}

func serveProxyStream(w http.ResponseWriter, r *http.Request, name string, modtime time.Time, rc io.ReadCloser) {
	defer rc.Close()
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, rc); err != nil {
		httpx.WriteError(w, err)
		return
	}
	http.ServeContent(w, r, name, modtime, bytes.NewReader(buf.Bytes()))
}

func (h *Handler) DownloadObject(w http.ResponseWriter, r *http.Request, fileId server.FileId) {
	p := principalOf(r)
	res, err := h.Objects.Download(r.Context(), p.TenantID, string(fileId), func(action string) {
		_ = h.Audit.Write(r.Context(), &domain.AuditEvent{
			TenantID: p.TenantID, Action: action, PrincipalType: string(p.Type),
			PrincipalID: p.ID, TargetType: "object", TargetID: string(fileId),
		})
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	if res.RedirectURL != "" {
		http.Redirect(w, r, res.RedirectURL, http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", res.ContentType)
	w.Header().Set("Content-Disposition", "attachment; filename=\""+strings.ReplaceAll(res.Object.Name, "\"", "")+"\"")
	serveProxyStream(w, r, res.Object.Name, res.Object.UpdatedAt, res.Stream)
}

// DownloadPublicObject serves an object with no authentication. It is only
// reachable when the owner set visibility=public — the usecase rejects
// private files. Content is served inline (so <img> works) but with
// X-Content-Type-Options: nosniff and a content-type allowlist so an uploaded
// .html/.svg/.php can't execute in the caller's browser.
func (h *Handler) DownloadPublicObject(w http.ResponseWriter, r *http.Request, fileId string) {
	res, err := h.Objects.PublicDownload(r.Context(), fileId)
	if err != nil {
		if httpx.IsNotFound(err) {
			w.WriteHeader(http.StatusNotFound)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte("<!doctype html><html><body><h1>Object not found or private</h1></body></html>"))
			return
		}
		httpx.WriteError(w, err)
		return
	}
	if res.RedirectURL != "" {
		http.Redirect(w, r, res.RedirectURL, http.StatusFound)
		return
	}
	setPublicContentHeaders(w, res)
	serveProxyStream(w, r, res.Object.Name, res.Object.UpdatedAt, res.Stream)
}

// setPublicContentHeaders applies safe headers for public (unauth) content.
// Inline only for display-safe media types; everything else forces download
// as application/octet-stream so uploaded HTML/SVG/JS can never run.
func setPublicContentHeaders(w http.ResponseWriter, res *object.DownloadResult) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Security-Policy", "default-src 'none'; sandbox")
	if isInlineSafeType(res.ContentType) {
		w.Header().Set("Content-Type", res.ContentType)
		w.Header().Set("Content-Disposition", "inline; filename=\""+strings.ReplaceAll(res.Object.Name, "\"", "")+"\"")
	} else {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+strings.ReplaceAll(res.Object.Name, "\"", "")+"\"")
	}
}

// isInlineSafeType reports whether a content type can be rendered inline.
// Images, audio and video are safe to preview; anything that can execute
// (HTML, SVG-as-text, JS, XML, fonts, etc.) is forced to download.
func isInlineSafeType(ct string) bool {
	switch {
	case ct == "":
		return false
	case strings.HasPrefix(ct, "image/"):
		// allow png/jpeg/gif/webp/avif/bmp; svg can carry scripts → download
		if strings.Contains(ct, "svg") {
			return false
		}
		return true
	case strings.HasPrefix(ct, "audio/"):
		return true
	case strings.HasPrefix(ct, "video/"):
		return true
	}
	return false
}

// --- shares ---

func (h *Handler) ListShareLinks(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	links, err := h.Shares.ListByTenant(r.Context(), p.TenantID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, links)
}

func (h *Handler) CreateShareLink(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	var req struct {
		ObjectID string `json:"object_id"`
		TTL      int    `json:"ttl"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	l, err := h.Shares.CreateSigned(r.Context(), p.TenantID, req.ObjectID, p.ID, req.TTL)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, l)
}

func (h *Handler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	var req struct {
		ObjectID string `json:"object_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	l, err := h.Shares.CreateShort(r.Context(), p.TenantID, req.ObjectID, p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, l)
}

func (h *Handler) RevokeShareLink(w http.ResponseWriter, r *http.Request, linkId server.LinkId) {
	p := principalOf(r)
	if err := h.Shares.Revoke(r.Context(), p.TenantID, string(linkId)); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ShareLinkStats(w http.ResponseWriter, r *http.Request, linkId server.LinkId) {
	p := principalOf(r)
	objs, err := h.Shares.ListByObject(r.Context(), p.TenantID, "")
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	_ = objs
	data(w, http.StatusOK, map[string]interface{}{"link_id": string(linkId)})
}

func (h *Handler) ResolveShortLink(w http.ResponseWriter, r *http.Request, slug string) {
	obj, _, err := h.Shares.Resolve(r.Context(), slug)
	if err != nil {
		// HTML 410 page, not the JSON envelope (§3.3)
		w.WriteHeader(http.StatusGone)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte("<!doctype html><html><body><h1>This link has expired</h1></body></html>"))
		return
	}
	res, err := h.Objects.Download(r.Context(), obj.TenantID, obj.ID, nil)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	if res.RedirectURL != "" {
		http.Redirect(w, r, res.RedirectURL, http.StatusFound)
		return
	}
	w.Header().Set("Content-Type", res.ContentType)
	serveProxyStream(w, r, res.Object.Name, res.Object.UpdatedAt, res.Stream)
}

// --- applications & keys ---

func (h *Handler) CreateTenantKey(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		TenantID       string   `json:"tenant_id"`
		Name           string   `json:"name"`
		ScopeFolderIDs []string `json:"scope_folder_ids"`
		Permissions    []string `json:"permissions"`
		ExpiresAt      *string  `json:"expires_at"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	// Platform admin can create keys for any tenant; members only their own.
	targetTenant := p.TenantID
	if req.TenantID != "" {
		if !p.IsPlatformAdmin {
			httpx.Error(w, http.StatusForbidden, "forbidden")
			return
		}
		targetTenant = req.TenantID
	}
	var exp *time.Time
	if req.ExpiresAt != nil {
		if t, err := time.Parse(time.RFC3339, *req.ExpiresAt); err == nil {
			exp = &t
		}
	}
	key, err := h.APIKeys.CreateTenantKey(r.Context(), targetTenant, req.Name, req.ScopeFolderIDs, req.Permissions, exp)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, key)
}

func (h *Handler) ListAllKeys(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	keys, err := h.APIKeys.ListKeysPage(r.Context(), p.TenantID, p.IsPlatformAdmin)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, keys)
}

func (h *Handler) RevokeKeyAny(w http.ResponseWriter, r *http.Request, keyId server.KeyId) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	if err := h.APIKeys.RevokeKeyAny(r.Context(), p.TenantID, string(keyId)); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListApplications(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	apps, err := h.APIKeys.List(r.Context(), p.TenantID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, apps)
}

func (h *Handler) RegisterApplication(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	a, err := h.APIKeys.Register(r.Context(), p.TenantID, req.Name, req.Description)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, a)
}

func (h *Handler) DeleteApplication(w http.ResponseWriter, r *http.Request, applicationId server.ApplicationId) {
	p := principalOf(r)
	if err := h.APIKeys.Delete(r.Context(), p.TenantID, string(applicationId)); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ListAccessKeys(w http.ResponseWriter, r *http.Request, applicationId server.ApplicationId) {
	p := principalOf(r)
	keys, err := h.APIKeys.ListKeys(r.Context(), p.TenantID, string(applicationId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, keys)
}

func (h *Handler) CreateAccessKey(w http.ResponseWriter, r *http.Request, applicationId server.ApplicationId) {
	p := principalOf(r)
	var req struct {
		Name           string   `json:"name"`
		ScopeFolderIDs []string `json:"scope_folder_ids"`
		Permissions    []string `json:"permissions"`
		ExpiresAt      *string  `json:"expires_at"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	var exp *time.Time
	if req.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err == nil {
			exp = &t
		}
	}
	key, err := h.APIKeys.CreateKey(r.Context(), p.TenantID, string(applicationId), req.Name, req.ScopeFolderIDs, req.Permissions, exp)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, key)
}

func (h *Handler) RevokeAccessKey(w http.ResponseWriter, r *http.Request, applicationId server.ApplicationId, keyId server.KeyId) {
	p := principalOf(r)
	if err := h.APIKeys.RevokeKey(r.Context(), p.TenantID, string(keyId)); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- grants ---

func (h *Handler) CreateGrant(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	var req struct {
		FolderID      string   `json:"folder_id"`
		PrincipalType string   `json:"principal_type"`
		PrincipalID   string   `json:"principal_id"`
		Permissions   []string `json:"permissions"`
		ExpiresAt     *string  `json:"expires_at"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	var exp *time.Time
	if req.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err == nil {
			exp = &t
		}
	}
	g, err := h.Grants.Create(r.Context(), p.TenantID, req.FolderID, req.PrincipalType, req.PrincipalID, req.Permissions, exp, p.ID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, g)
}

func (h *Handler) ListFolderGrants(w http.ResponseWriter, r *http.Request, folderId server.FolderId) {
	p := principalOf(r)
	gs, err := h.Grants.ListByFolder(r.Context(), p.TenantID, string(folderId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, gs)
}

func (h *Handler) RevokeGrant(w http.ResponseWriter, r *http.Request, grantId server.GrantId) {
	p := principalOf(r)
	if err := h.Grants.Revoke(r.Context(), p.TenantID, string(grantId)); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- archives & jobs ---

func (h *Handler) ExtractArchive(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	var req struct {
		FileID         string `json:"file_id"`
		TargetFolderID string `json:"target_folder_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	jobID, err := h.Jobs.Enqueue(r.Context(), p.TenantID, "extract", map[string]interface{}{
		"file_id": req.FileID, "target_folder_id": req.TargetFolderID,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusAccepted, map[string]interface{}{"job_id": jobID})
}

func (h *Handler) CreateBundle(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	var req struct {
		FolderID string `json:"folder_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	jobID, err := h.Jobs.Enqueue(r.Context(), p.TenantID, "bundle", map[string]interface{}{"folder_id": req.FolderID})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusAccepted, map[string]interface{}{"job_id": jobID})
}

func (h *Handler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	var req struct {
		SourceStorageID string `json:"source_storage_id"`
		TargetStorageID string `json:"target_storage_id"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	if req.SourceStorageID == "" || req.TargetStorageID == "" || req.SourceStorageID == req.TargetStorageID {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	jobID, err := h.Jobs.Enqueue(r.Context(), p.TenantID, "transfer", map[string]interface{}{
		"source_storage_id": req.SourceStorageID,
		"target_storage_id": req.TargetStorageID,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusAccepted, map[string]interface{}{"job_id": jobID})
}

func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	jobs, err := h.Jobs.List(r.Context(), p.TenantID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, jobs)
}

func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request, jobId server.JobId) {
	p := principalOf(r)
	j, err := h.Jobs.Get(r.Context(), p.TenantID, string(jobId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, j)
}

// --- audit ---

func (h *Handler) ListAuditEvents(w http.ResponseWriter, r *http.Request, params server.ListAuditEventsParams) {
	p := principalOf(r)
	limit := 50
	if params.Limit != nil {
		limit = *params.Limit
	}
	targetType, targetID, action := "", "", ""
	if params.TargetType != nil {
		targetType = *params.TargetType
	}
	if params.TargetId != nil {
		targetID = *params.TargetId
	}
	if params.Action != nil {
		action = *params.Action
	}
	evs, err := h.Audit.Query(r.Context(), p.TenantID, audit.ListFilter{
		TargetType: targetType, TargetID: targetID, Action: action, Limit: limit,
	})
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, evs)
}

// --- usage ---

func (h *Handler) MyUsage(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	// Live numbers from the tenant record (denormalized counters) so the usage
	// page works before the hourly metering ticker has produced snapshots.
	live := &domain.UsageSnapshot{}
	if t, err := h.Tenants.Get(r.Context(), p.TenantID); err == nil {
		live.BytesStored = t.UsedBytes
		live.ObjectCount = t.UsedObjects
		live.Period = "live"
	}
	if snap, err := h.Usage.Latest(r.Context(), p.TenantID); err == nil {
		// merge: prefer live counters for bytes/objects, snapshot for egress
		live.EgressBytes = snap.EgressBytes
		live.RequestCount = snap.RequestCount
	}
	data(w, http.StatusOK, live)
}

func (h *Handler) EstimatedCost(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	est, err := h.Usage.EstimatedCost(r.Context(), p.TenantID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, est)
}

// --- admin ---

func (h *Handler) ListAvailableBackends(w http.ResponseWriter, r *http.Request) {
	p := principalOf(r)
	if p == nil {
		httpx.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	bs, err := h.Admin.ListAvailable(r.Context(), p.TenantID)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, bs)
}

func (h *Handler) ListBackends(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	bs, err := h.Admin.ListBackends(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, bs)
}

func (h *Handler) CreateBackend(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	var req struct {
		Name        string                 `json:"name"`
		Driver      string                 `json:"driver"`
		Config      map[string]interface{} `json:"config"`
		Credentials map[string]interface{} `json:"credentials"`
		RateCard    map[string]interface{} `json:"rate_card"`
	}
	if err := decodeBody(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request")
		return
	}
	b, err := h.Admin.CreateBackend(r.Context(), req.Name, req.Driver, req.Config, req.Credentials, req.RateCard)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusCreated, b)
}

func (h *Handler) GetBackend(w http.ResponseWriter, r *http.Request, backendId server.BackendId) {
	if !requireAdmin(w, r) {
		return
	}
	b, err := h.Admin.GetBackend(r.Context(), string(backendId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, b)
}

func (h *Handler) DeleteBackend(w http.ResponseWriter, r *http.Request, backendId server.BackendId) {
	if !requireAdmin(w, r) {
		return
	}
	if err := h.Admin.DeleteBackend(r.Context(), string(backendId)); err != nil {
		httpx.WriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CheckBackendHealth(w http.ResponseWriter, r *http.Request, backendId server.BackendId) {
	if !requireAdmin(w, r) {
		return
	}
	b, err := h.Admin.CheckHealth(r.Context(), string(backendId))
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, b)
}

func (h *Handler) ListAllTenants(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	ts, err := h.Admin.ListAllTenants(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, ts)
}

func (h *Handler) AdminUsage(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	rows, err := h.Admin.TenantUsage(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, rows)
}

func (h *Handler) InstallStats(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	stats, err := h.Admin.InstallStats(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}
	data(w, http.StatusOK, stats)
}

func strPtr(s string) *string { return &s }

// helpers to satisfy unused-imports guard in some builds
