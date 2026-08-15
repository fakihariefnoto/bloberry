# Task group — 19 page: members

**Depends on:** `01-setup`, `02-design-tokens`, `03-routing`, `05-core-infra`, `06-shared-components` (DataTable, StatusPill, FormField, ConfirmDestructive, RelativeTime). **Blocks:** `33-flows` (invite flow). **Mockup:** [`web/mockup/members.md`](../mockup/members.md).

- [ ] **Layout — desktop** per the mockup: PageHeader "Members" with `Invite`, DataTable (member, email, role, joined, `⋮`), the Pending invitations panel below with Resend, the role-floor footer caption.
- [ ] **Layout — mobile**: stacked cards.
- [ ] **`invite-member` modal** — email + role select; on send the invite appears in the pending list. Only admins/owners can invite.
- [ ] **Pending invites** — show email, role, expiry; **Resend** regenerates the invite (replacing the old token — a previously-emailed link stops working). Expired invites drop off and self-delete server-side (`ERD.md` invitations TTL).
- [ ] **Role change** — inline select; **demoting the owner requires typed-name confirmation** (the tenant's sole authority changing hands); any other role change confirms once.
- [ ] **Remove member** — `ConfirmDestructive`; removing the **last owner** is refused ("A tenant must have at least one owner."); removing a member with grants keeps the grants inert but retained for the audit trail (`ERD.md` grants).
- [ ] **The owner row** — no role select, no remove (UI lock, not hidden authority).
- [ ] **Role floor caption** — "Roles set a floor. Folder grants add access on top — a viewer can still be granted write on one folder." (PRD D7, no deny rules to hunt for).
- [ ] **Empty states** — no members → "Invite your team · People join through emailed invitations"; pending-empty → no panel at all.
- [ ] **A11y** — role select keyboard-accessible with visible focus; confirm dialogs trap focus and return it to the triggering row.

**tests:** owner cannot be demoted/removed; last-owner removal refused; demote-owner typed-name; Resend invalidates the old link; role-floor caption present.
