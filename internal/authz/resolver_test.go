package authz

import "testing"

func TestResolve(t *testing.T) {
	tests := []struct {
		name      string
		p         Principal
		action    Permission
		folderID  string
		ancestors []string
		want      Decision
	}{
		{"platform admin always allow", Principal{IsPlatformAdmin: true}, PermRead, "f1", nil, Allow},
		{"no tenant deny", Principal{}, PermRead, "f1", nil, Deny},
		{"owner reads", Principal{TenantID: "t", Role: RoleTenantOwner}, PermRead, "f1", nil, Allow},
		{"owner writes", Principal{TenantID: "t", Role: RoleTenantOwner}, PermWrite, "f1", nil, Allow},
		{"owner deletes", Principal{TenantID: "t", Role: RoleTenantOwner}, PermDelete, "f1", nil, Allow},
		{"owner shares", Principal{TenantID: "t", Role: RoleTenantOwner}, PermShare, "f1", nil, Allow},
		{"owner admin", Principal{TenantID: "t", Role: RoleTenantOwner}, PermAdmin, "f1", nil, Allow},
		{"tenant_admin reads", Principal{TenantID: "t", Role: RoleTenantAdmin}, PermRead, "f1", nil, Allow},
		{"tenant_admin admin", Principal{TenantID: "t", Role: RoleTenantAdmin}, PermAdmin, "f1", nil, Allow},
		{"member reads", Principal{TenantID: "t", Role: RoleMember}, PermRead, "f1", nil, Allow},
		{"member writes", Principal{TenantID: "t", Role: RoleMember}, PermWrite, "f1", nil, Allow},
		{"member shares", Principal{TenantID: "t", Role: RoleMember}, PermShare, "f1", nil, Allow},
		{"member cannot delete", Principal{TenantID: "t", Role: RoleMember}, PermDelete, "f1", nil, Deny},
		{"member cannot admin", Principal{TenantID: "t", Role: RoleMember}, PermAdmin, "f1", nil, Deny},
		{"viewer reads", Principal{TenantID: "t", Role: RoleViewer}, PermRead, "f1", nil, Allow},
		{"viewer cannot write", Principal{TenantID: "t", Role: RoleViewer}, PermWrite, "f1", nil, Deny},
		{"viewer cannot share", Principal{TenantID: "t", Role: RoleViewer}, PermShare, "f1", nil, Deny},
		{"viewer cannot delete", Principal{TenantID: "t", Role: RoleViewer}, PermDelete, "f1", nil, Deny},
		{"viewer cannot admin", Principal{TenantID: "t", Role: RoleViewer}, PermAdmin, "f1", nil, Deny},

		{"key scope own folder", Principal{TenantID: "t", Role: RoleMember, ScopeFolderIDs: []string{"f1"}}, PermRead, "f1", nil, Allow},
		{"key scope descendant", Principal{TenantID: "t", Role: RoleMember, ScopeFolderIDs: []string{"f1"}}, PermRead, "f2", []string{"f1"}, Allow},
		{"key scope ancestor target denied", Principal{TenantID: "t", Role: RoleMember, ScopeFolderIDs: []string{"f2"}}, PermRead, "f1", nil, Deny},
		{"key scope sibling denied", Principal{TenantID: "t", Role: RoleMember, ScopeFolderIDs: []string{"f1"}}, PermRead, "f3", []string{"root"}, Deny},
		{"key scope empty means whole tenant", Principal{TenantID: "t", Role: RoleMember}, PermRead, "f9", nil, Allow},

		{"key permissions narrow", Principal{TenantID: "t", Role: RoleMember, KeyPermissions: []Permission{PermRead}}, PermRead, "f1", nil, Allow},
		{"key permissions deny write", Principal{TenantID: "t", Role: RoleMember, KeyPermissions: []Permission{PermRead}}, PermWrite, "f1", nil, Deny},
		{"key permissions empty no narrowing", Principal{TenantID: "t", Role: RoleMember}, PermRead, "f1", nil, Allow},

		{"grant on target folder", Principal{TenantID: "t", Role: RoleViewer, Grants: []Grant{{FolderID: "f1", Perms: []Permission{PermWrite}}}}, PermWrite, "f1", nil, Allow},
		{"grant on ancestor", Principal{TenantID: "t", Role: RoleViewer, Grants: []Grant{{FolderID: "root", Perms: []Permission{PermWrite}}}}, PermWrite, "f2", []string{"root"}, Allow},
		{"deeper grant wins", Principal{TenantID: "t", Role: RoleViewer, Grants: []Grant{{FolderID: "root", Perms: []Permission{PermWrite}}, {FolderID: "f1", Perms: []Permission{PermRead}}}}, PermRead, "f2", []string{"root", "f1"}, Allow},
		{"two grants same depth first kept", Principal{TenantID: "t", Role: RoleViewer, Grants: []Grant{{FolderID: "f1", Perms: []Permission{PermRead}}, {FolderID: "f1", Perms: []Permission{PermWrite}}}}, PermWrite, "f1", nil, Deny},
		{"scoped key grant denied by scope", Principal{TenantID: "t", Role: RoleMember, ScopeFolderIDs: []string{"f2"}, Grants: []Grant{{FolderID: "f1", Perms: []Permission{PermRead}}}}, PermRead, "f1", nil, Deny},
		{"expired grant ignored", Principal{TenantID: "t", Role: RoleViewer, Grants: []Grant{{FolderID: "f1", Perms: []Permission{PermWrite}, Expired: true}}}, PermWrite, "f1", nil, Deny},
		{"revoked grant ignored", Principal{TenantID: "t", Role: RoleViewer, Grants: []Grant{{FolderID: "f1", Perms: []Permission{PermWrite}, Revoked: true}}}, PermWrite, "f1", nil, Deny},
		{"grant not on ancestor denied", Principal{TenantID: "t", Role: RoleViewer, Grants: []Grant{{FolderID: "other", Perms: []Permission{PermWrite}}}}, PermWrite, "f1", nil, Deny},
		{"viewer grant read adds", Principal{TenantID: "t", Role: RoleViewer, Grants: []Grant{{FolderID: "f1", Perms: []Permission{PermRead}}}}, PermRead, "f1", nil, Allow},

		{"app key whole tenant read", Principal{TenantID: "t", Type: PrincipalApplication, Role: RoleMember}, PermRead, "f1", nil, Allow},
		{"deny by default", Principal{TenantID: "t", Role: ""}, PermRead, "f1", nil, Deny},
		{"unknown action denied", Principal{TenantID: "t", Role: RoleTenantOwner}, Permission("explode"), "f1", nil, Deny},
		{"unknown role denied", Principal{TenantID: "t", Role: Role("superuser")}, PermRead, "f1", nil, Deny},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Resolve(tt.p, tt.action, tt.folderID, tt.ancestors)
			if got != tt.want {
				t.Errorf("Resolve() = %v, want %v", got, tt.want)
			}
		})
	}
}
