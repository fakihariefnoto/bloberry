package domain

import "time"

type User struct {
	ID              string        `bson:"_id" json:"id"`
	Email           string        `bson:"email" json:"email"`
	PasswordHash    *string       `bson:"password_hash,omitempty" json:"-"`
	DisplayName     string        `bson:"display_name" json:"display_name"`
	PlatformRole    *string       `bson:"platform_role,omitempty" json:"platform_role,omitempty"`
	EmailVerified   bool          `bson:"email_verified" json:"email_verified"`
	LastLoginAt     *time.Time    `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
	OAuthIdentities []OAuthID     `bson:"oauth_identities,omitempty" json:"oauth_identities,omitempty"`
	Settings        UserSettings  `bson:"settings" json:"settings"`
	TOTP            *TOTPConfig   `bson:"totp,omitempty" json:"totp,omitempty"`
	CreatedAt       time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time     `bson:"updated_at" json:"updated_at"`
}

type OAuthID struct {
	Provider        string    `bson:"provider" json:"provider"`
	ProviderUserID  string    `bson:"provider_user_id" json:"provider_user_id"`
	EmailAtLink     string    `bson:"email_at_link" json:"email_at_link"`
	CreatedAt       time.Time `bson:"created_at" json:"created_at"`
}

type UserSettings struct {
	NotificationsEnabled bool   `bson:"notifications_enabled" json:"notifications_enabled"`
	BiometricUnlockEnabled bool `bson:"biometric_unlock_enabled" json:"biometric_unlock_enabled"`
	Locale               string `bson:"locale" json:"locale"`
	DefaultTenantID      string `bson:"default_tenant_id" json:"default_tenant_id"`
}

type BackupCode struct {
	Hash string `bson:"hash" json:"hash"`
	Used bool   `bson:"used" json:"used"`
}

type TOTPConfig struct {
	SecretEncrypted string       `bson:"secret_encrypted" json:"-"`
	Enabled         bool         `bson:"enabled" json:"enabled"`
	BackupCodes     []BackupCode `bson:"backup_codes" json:"-"`
	EnabledAt       *time.Time   `bson:"enabled_at,omitempty" json:"enabled_at,omitempty"`
}

type Tenant struct {
	ID               string     `bson:"_id" json:"id"`
	Name             string     `bson:"name" json:"name"`
	Slug             string     `bson:"slug" json:"slug"`
	DefaultBackendID string     `bson:"default_backend_id,omitempty" json:"default_storage_id,omitempty"`
	// StorageEngines is the set of engines this tenant may use (in addition to
	// install-level ones). Assigning an engine here makes it selectable in the
	// Files storage switcher for this tenant.
	StorageEngines    []string   `bson:"storage_engines,omitempty" json:"storage_engines,omitempty"`
	QuotaBytes       int64      `bson:"quota_bytes" json:"quota_bytes"`
	QuotaObjects     int64      `bson:"quota_objects" json:"quota_objects"`
	UsedBytes        int64      `bson:"used_bytes" json:"used_bytes"`
	UsedObjects      int64      `bson:"used_objects" json:"used_objects"`
	Status           string     `bson:"status" json:"status"`
	Billing          map[string]interface{} `bson:"billing,omitempty" json:"billing,omitempty"`
	CreatedAt        time.Time  `bson:"created_at" json:"created_at"`
}

type Membership struct {
	ID        string    `bson:"_id" json:"id"`
	UserID    string    `bson:"user_id" json:"user_id"`
	TenantID  string    `bson:"tenant_id" json:"tenant_id"`
	Role      string    `bson:"role" json:"role"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

type Invitation struct {
	ID         string     `bson:"_id" json:"id"`
	TenantID   string     `bson:"tenant_id" json:"tenant_id"`
	Email      string     `bson:"email" json:"email"`
	Role       string     `bson:"role" json:"role"`
	TokenHash  string     `bson:"token_hash" json:"-"`
	InvitedBy  string     `bson:"invited_by" json:"invited_by"`
	ExpiresAt  time.Time  `bson:"expires_at" json:"expires_at"`
	AcceptedAt *time.Time `bson:"accepted_at,omitempty" json:"accepted_at,omitempty"`
	CreatedAt  time.Time  `bson:"created_at" json:"created_at"`
}

type Application struct {
	ID          string    `bson:"_id" json:"id"`
	TenantID    string    `bson:"tenant_id" json:"tenant_id"`
	Name        string    `bson:"name" json:"name"`
	Description string    `bson:"description" json:"description"`
	CreatedAt   time.Time `bson:"created_at" json:"created_at"`
}

type AccessKey struct {
	ID             string     `bson:"_id" json:"id"`
	TenantID       string     `bson:"tenant_id" json:"tenant_id"`
	Name           string     `bson:"name" json:"name"`
	ApplicationID  *string    `bson:"application_id,omitempty" json:"application_id,omitempty"`
	UserID         *string    `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Prefix         string     `bson:"prefix" json:"prefix"`
	SecretHash     string     `bson:"secret_hash" json:"-"`
	LastFour       string     `bson:"last_four" json:"last_four"`
	ScopeFolderIDs []string   `bson:"scope_folder_ids" json:"scope_folder_ids"`
	Permissions    []string   `bson:"permissions" json:"permissions"`
	ExpiresAt      *time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	LastUsedAt     *time.Time `bson:"last_used_at,omitempty" json:"last_used_at,omitempty"`
	LastUsedIP     string     `bson:"last_used_ip" json:"last_used_ip"`
	RevokedAt      *time.Time `bson:"revoked_at,omitempty" json:"revoked_at,omitempty"`
	CreatedAt      time.Time  `bson:"created_at" json:"created_at"`
}

type StorageBackend struct {
	ID                 string                 `bson:"_id" json:"id"`
	TenantID           *string                `bson:"tenant_id,omitempty" json:"tenant_id,omitempty"`
	Driver             string                 `bson:"driver" json:"driver"`
	Name               string                 `bson:"name" json:"name"`
	Config             map[string]interface{} `bson:"config" json:"config"`
	CredentialsEncrypted []byte               `bson:"credentials_encrypted" json:"-"`
	Capabilities       map[string]interface{} `bson:"capabilities" json:"capabilities"`
	RateCard           map[string]interface{} `bson:"rate_card,omitempty" json:"rate_card,omitempty"`
	HealthStatus       string                 `bson:"health_status" json:"health_status"`
	HealthError        string                 `bson:"health_error,omitempty" json:"health_error,omitempty"`
	HealthCheckedAt    *time.Time             `bson:"health_checked_at,omitempty" json:"health_checked_at,omitempty"`
	CreatedAt          time.Time              `bson:"created_at" json:"created_at"`
}

type Folder struct {
	ID        string    `bson:"_id" json:"id"`
	TenantID  string    `bson:"tenant_id" json:"tenant_id"`
	ParentID  *string   `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Path      string    `bson:"path" json:"path"`
	Ancestors []string  `bson:"ancestors" json:"ancestors"`
	Depth     int       `bson:"depth" json:"depth"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type Object struct {
	ID          string     `bson:"_id" json:"id"`
	TenantID    string     `bson:"tenant_id" json:"tenant_id"`
	FolderID    string     `bson:"folder_id" json:"folder_id"`
	Ancestors   []string   `bson:"ancestors" json:"ancestors"`
	Name        string     `bson:"name" json:"name"`
	BackendID   string     `bson:"backend_id" json:"storage_id"`
	StorageKey  string     `bson:"storage_key" json:"storage_key"`
	State       string     `bson:"state" json:"state"`
	SizeBytes   int64      `bson:"size_bytes" json:"size_bytes"`
	ContentType string     `bson:"content_type" json:"content_type"`
	ContentHash string     `bson:"content_hash" json:"content_hash"`
	Visibility  string     `bson:"visibility" json:"visibility"`
	UploadedBy  string     `bson:"uploaded_by" json:"uploaded_by"`
	DeletedAt   *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
	CreatedAt   time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `bson:"updated_at" json:"updated_at"`
}

type MultipartUpload struct {
	ID               string    `bson:"_id" json:"id"`
	ObjectID         string    `bson:"object_id" json:"object_id"`
	TenantID         string    `bson:"tenant_id" json:"tenant_id"`
	ProviderUploadID string    `bson:"provider_upload_id" json:"provider_upload_id"`
	PartSizeBytes    int64     `bson:"part_size_bytes" json:"part_size_bytes"`
	PartsReceived    []PartRec `bson:"parts_received" json:"parts_received"`
	ExpiresAt        time.Time `bson:"expires_at" json:"expires_at"`
	CreatedAt        time.Time `bson:"created_at" json:"created_at"`
}

type PartRec struct {
	Part int    `bson:"part" json:"part"`
	ETag string `bson:"etag" json:"etag"`
}

type Grant struct {
	ID           string     `bson:"_id" json:"id"`
	TenantID     string     `bson:"tenant_id" json:"tenant_id"`
	FolderID     string     `bson:"folder_id" json:"folder_id"`
	PrincipalType string    `bson:"principal_type" json:"principal_type"`
	PrincipalID  string     `bson:"principal_id" json:"principal_id"`
	Permissions  []string   `bson:"permissions" json:"permissions"`
	ExpiresAt    *time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	GrantedBy    string     `bson:"granted_by" json:"granted_by"`
	RevokedAt    *time.Time `bson:"revoked_at,omitempty" json:"revoked_at,omitempty"`
	CreatedAt    time.Time  `bson:"created_at" json:"created_at"`
}

type ShareLink struct {
	ID            string     `bson:"_id" json:"id"`
	TenantID      string     `bson:"tenant_id" json:"tenant_id"`
	ObjectID      string     `bson:"object_id" json:"object_id"`
	Kind          string     `bson:"kind" json:"kind"`
	Slug          string     `bson:"slug,omitempty" json:"slug,omitempty"`
	TokenHash     string     `bson:"token_hash,omitempty" json:"-"`
	ExpiresAt     *time.Time `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
	HitCount      int64      `bson:"hit_count" json:"hit_count"`
	LastAccessedAt *time.Time `bson:"last_accessed_at,omitempty" json:"last_accessed_at,omitempty"`
	CreatedBy     string     `bson:"created_by" json:"created_by"`
	RevokedAt     *time.Time `bson:"revoked_at,omitempty" json:"revoked_at,omitempty"`
	CreatedAt     time.Time  `bson:"created_at" json:"created_at"`
}

type Job struct {
	ID             string    `bson:"_id" json:"id"`
	TenantID       string    `bson:"tenant_id" json:"tenant_id"`
	Kind           string    `bson:"kind" json:"kind"`
	State          string    `bson:"state" json:"state"`
	Payload        map[string]interface{} `bson:"payload" json:"payload"`
	Result         map[string]interface{} `bson:"result,omitempty" json:"result,omitempty"`
	ProgressDone   int       `bson:"progress_done" json:"progress_done"`
	ProgressTotal  int       `bson:"progress_total" json:"progress_total"`
	FailureCode    string    `bson:"failure_code,omitempty" json:"failure_code,omitempty"`
	FailureMessage string    `bson:"failure_message,omitempty" json:"failure_message,omitempty"`
	Attempts       int       `bson:"attempts" json:"attempts"`
	StartedAt      *time.Time `bson:"started_at,omitempty" json:"started_at,omitempty"`
	FinishedAt     *time.Time `bson:"finished_at,omitempty" json:"finished_at,omitempty"`
	CreatedAt      time.Time `bson:"created_at" json:"created_at"`
}

type UsageSnapshot struct {
	ID           string  `bson:"_id" json:"id"`
	TenantID     string  `bson:"tenant_id" json:"tenant_id"`
	Period       string  `bson:"period" json:"period"`
	BytesStored  int64   `bson:"bytes_stored" json:"bytes_stored"`
	ObjectCount  int64   `bson:"object_count" json:"object_count"`
	EgressBytes  int64   `bson:"egress_bytes" json:"egress_bytes"`
	RequestCount int64   `bson:"request_count" json:"request_count"`
	EstimatedCost float64 `bson:"estimated_cost" json:"estimated_cost"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

type AuditEvent struct {
	ID           string    `bson:"_id" json:"id"`
	TenantID     string    `bson:"tenant_id" json:"tenant_id"`
	Action       string    `bson:"action" json:"action"`
	PrincipalType string   `bson:"principal_type" json:"principal_type"`
	PrincipalID  string    `bson:"principal_id" json:"principal_id"`
	TargetType   string    `bson:"target_type,omitempty" json:"target_type,omitempty"`
	TargetID     string    `bson:"target_id,omitempty" json:"target_id,omitempty"`
	Metadata     map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	IP           string    `bson:"ip" json:"ip"`
	UserAgent    string    `bson:"user_agent" json:"user_agent"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}
