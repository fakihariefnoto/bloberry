package httpx

import (
	"errors"
	"net/http"
)

type Code string

const (
	ErrInvalidCredentials Code = "invalid_credentials"
	ErrRefreshInvalid     Code = "refresh_invalid"
	ErrKeyRevoked         Code = "key_revoked"
	ErrKeyExpired         Code = "key_expired"
	ErrOtpInvalid         Code = "otp_invalid"
	ErrOtpRateLimited     Code = "otp_rate_limited"
	ErrOauthInvalid       Code = "oauth_invalid"
	ErrPairInvalid        Code = "pair_invalid"
	ErrConfigInvalid      Code = "config_invalid"
	ErrConfigExpired      Code = "config_expired"
	ErrTotpRequired       Code = "totp_required"
	ErrTotpInvalid        Code = "totp_invalid"
	ErrNoInvitation       Code = "no_invitation"
	ErrForbidden          Code = "forbidden"
	ErrQuotaExceeded      Code = "quota_exceeded"
	ErrNameConflict       Code = "name_conflict"
	ErrFolderCycle        Code = "folder_cycle"
	ErrObjectPending      Code = "object_pending"
	ErrBackendUnreachable Code = "backend_unreachable"
	ErrArchiveRejected    Code = "archive_rejected"
	ErrLinkExpired        Code = "link_expired"
	ErrPayloadTooLarge    Code = "payload_too_large"
	ErrInviteInvalid      Code = "invite_invalid"
	ErrResetTokenInvalid  Code = "reset_token_invalid"
	ErrRateLimited        Code = "rate_limited"
	ErrNotFound           Code = "not_found"
	ErrBadRequest         Code = "bad_request"
	ErrUnauthorized       Code = "unauthorized"
	ErrConflict           Code = "conflict"
	ErrInternal           Code = "internal"
)

type HTTPError struct {
	Code    Code
	Status  int
	Content string
}

func (e *HTTPError) Error() string {
	return string(e.Code)
}

func NewError(code Code, status int) *HTTPError {
	return &HTTPError{Code: code, Status: status}
}

func NewErrorContent(code Code, status int, content string) *HTTPError {
	return &HTTPError{Code: code, Status: status, Content: content}
}

// ErrResourceNotFound is the standard not-found sentinel repositories return.
var ErrResourceNotFound = NewError(ErrNotFound, http.StatusNotFound)

func IsNotFound(err error) bool {
	var he *HTTPError
	return errors.As(err, &he) && he.Code == ErrNotFound
}

func From(err error) *HTTPError {
	if err == nil {
		return nil
	}
	var he *HTTPError
	if errors.As(err, &he) {
		return he
	}
	return NewError(ErrInternal, http.StatusInternalServerError)
}

func WriteError(w http.ResponseWriter, err error) {
	he := From(err)
	if he.Status == http.StatusInternalServerError {
		ErrorWithContent(w, he.Status, string(he.Code), "internal error")
		return
	}
	if he.Content != "" {
		ErrorWithContent(w, he.Status, string(he.Code), he.Content)
		return
	}
	Error(w, he.Status, string(he.Code))
}
