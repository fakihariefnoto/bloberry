package authz

type PrincipalType string

const (
	PrincipalUser        PrincipalType = "user"
	PrincipalApplication PrincipalType = "application"
)

type Role string

const (
	RoleTenantOwner  Role = "tenant_owner"
	RoleTenantAdmin  Role = "tenant_admin"
	RoleMember       Role = "member"
	RoleViewer       Role = "viewer"
)

type Permission string

const (
	PermRead   Permission = "read"
	PermWrite  Permission = "write"
	PermDelete Permission = "delete"
	PermShare  Permission = "share"
	PermAdmin  Permission = "admin"
)

type Grant struct {
	FolderID string
	Expired  bool
	Revoked  bool
	Perms    []Permission
}

type Principal struct {
	Type            PrincipalType
	ID              string
	TenantID        string
	Role            Role
	IsPlatformAdmin bool
	ScopeFolderIDs  []string
	KeyPermissions  []Permission
	Grants          []Grant
}

type Decision int

const (
	Deny  Decision = 0
	Allow Decision = 1
)

func has(perms []Permission, want Permission) bool {
	for _, p := range perms {
		if p == want {
			return true
		}
	}
	return false
}

func roleFloor(role Role, action Permission) bool {
	switch action {
	case PermRead:
		return role == RoleTenantOwner || role == RoleTenantAdmin || role == RoleMember || role == RoleViewer
	case PermWrite, PermShare:
		return role == RoleTenantOwner || role == RoleTenantAdmin || role == RoleMember
	case PermDelete, PermAdmin:
		return role == RoleTenantOwner || role == RoleTenantAdmin
	}
	return false
}

// Resolve is the pure permission function. No I/O. ADR-6, PRD G5.
func Resolve(p Principal, action Permission, folderID string, ancestors []string) Decision {
	if p.IsPlatformAdmin {
		return Allow
	}
	if p.TenantID == "" {
		return Deny
	}
	// Key scope: narrows only. Non-empty scope must contain the folder or an ancestor.
	if len(p.ScopeFolderIDs) > 0 {
		inScope := contains(p.ScopeFolderIDs, folderID)
		for _, a := range ancestors {
			if contains(p.ScopeFolderIDs, a) {
				inScope = true
				break
			}
		}
		if !inScope {
			return Deny
		}
	}
	// Key permissions: narrows only.
	if len(p.KeyPermissions) > 0 && !has(p.KeyPermissions, action) {
		return Deny
	}
	// Role floor.
	if roleFloor(p.Role, action) {
		return Allow
	}
	// Grants: deepest matching ancestor. Allow-only, adds.
	if g := deepestGrant(p, folderID, ancestors); g != nil && has(g.Perms, action) {
		return Allow
	}
	return Deny
}

func deepestGrant(p Principal, folderID string, ancestors []string) *Grant {
	var best *Grant
	bestDepth := -1
	consider := func(id string, depth int) {
		for i := range p.Grants {
			g := &p.Grants[i]
			if g.FolderID != id || g.Expired || g.Revoked {
				continue
			}
			if depth > bestDepth {
				best = g
				bestDepth = depth
			}
		}
	}
	// folderID is deepest: depth = len(ancestors). ancestors[i] is at depth i.
	consider(folderID, len(ancestors))
	for i := len(ancestors) - 1; i >= 0; i-- {
		consider(ancestors[i], i)
	}
	return best
}

func contains(xs []string, v string) bool {
	for _, x := range xs {
		if x == v {
			return true
		}
	}
	return false
}
