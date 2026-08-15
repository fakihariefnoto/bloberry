# Screen — more

## Purpose & context

- **User goal**: reach the occasional-use administrative surfaces — profile, settings, usage, applications — plus the tenant switcher. The "everything administrative" tab (`mobile/navigation.md`): deliberately a menu, not three more tabs.
- **Entry points**: tab tap.
- **Exit points**: Profile → `profile`; Settings → `settings`; Usage → `usage`; Applications → `applications`; tenant row → `tenant-switcher` sheet (switching lands on `files` at the new root); app version/about inline.
- **Data needed**: current user (name, avatar, role in this tenant), current tenant, storage usage summary (bytes/quota), and role-gated entries.

## States

- [x] Single state (static menu + user header + usage summary). No loading beyond the user header skeleton; nothing else fetches.

## Style reference

- **Components used**: settings-list pattern (`design-collection/mobile-screen/patterns.md`) — grouped rows with headers, `›` affordances, toggle/summary values right-aligned. User header (avatar, name, email, role badge) at top; tenant switcher at top with the user header.
- No token deltas.

## Wireframe — mobile

```
┌───────────────────────────┐
│  More                     │
├───────────────────────────┤
│  ┌──────────────────────┐ │
│  │ 👤 Jane Doe          │ │
│  │ jane@acme.dev        │ │
│  │ tenant_owner         │ │
│  └──────────────────────┘ │
│  ──────────────────────── │
│  TENANT                   │
│  Acme Inc              ▾  │
│  Storage  312/500 GB 62%  │
│  ──────────────────────── │
│  ACCOUNT                  │
│  Profile               ›  │
│  Settings              ›  │
│  ──────────────────────── │
│  ADMIN                    │
│  Usage                 ›  │
│  Applications          ›  │
│  ──────────────────────── │
│  ABOUT                    │
│  Version 0.1.0            │
│  Log out              ⏻   │
├───────────────────────────┤
│  Files  Uploads  Shares  ●│
└───────────────────────────┘
```

## Interactions

- **User header**: tap → `profile`.
- **Tenant row** → `tenant-switcher` sheet; switching **always lands on `files` at the new tenant's root** (`navigation.md` tenant-switch rule).
- **Usage row** shows the live quota summary inline (bytes/percent bar) — the admin glancing at "am I near quota" doesn't even need to open `usage`.
- **Admin group** (`Usage`, `Applications`) is shown only to `tenant_admin`+; a `member`/`viewer` sees the section header hidden entirely (not an empty group). The role badge in the header is the same gate.
- **Log out** sits at the bottom under ABOUT, visually separated (settings-list pattern). It's a plain action, not a nested confirm — logging out is reversible by logging back in. It does revoke the refresh token server-side.
- **Version** row is static text (`text.caption`), tappable only in debug builds (shows build commit).
- **A11y**: the user header's avatar is `role="presentation"`; the role badge has a text label; every row is ≥48dp.
