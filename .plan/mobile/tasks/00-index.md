# Build index — Mobile

Build-order table, dependency edges, and status for the Flutter mobile task list.

## Build order

| File | Covers | Status |
|---|---|---|
| `01-setup.md` | `flutter create`, app id `com.bloberry.app`, packages, folder tree | ☐ |
| `02-design-tokens.md` | ThemeData + token file from `design/style-guide.md` (every scale + component states) | ☐ |
| `03-routing.md` | go_router, 21 routes with typed args, 4-tab shell, 5 sheets, deep links, auth gates | ☐ |
| `04-core-infra.md` | Dio + envelope interceptor, snake_case JSON, token storage, validators, biometric unlock | ☐ |
| `05-screen-onboarding.md` … `25-screen-tenant-suspended.md` | One file per screen (21), grounded in its mockup | ☐ |
| `26-flows.md` | Cross-screen flows from `mobile/navigation.md` §Flow chains (5 flows + QR pair) | ☐ |
| `27-platform-config.md` | Permissions, icon/splash, background-upload services | ☐ |
| `28-testing.md` | Named high-risk-flow widget tests | ☐ |

## Dependency edges

- **`01-setup` blocks everything.**
- **`02-design-tokens` blocks every screen** (screens built before tokens exist get hardcoded one-off values).
- **`03-routing` blocks every screen** — screens render inside the shell; route args are typed per the route table.
- **`04-core-infra` blocks the API-touching screens** — `files`, `file-detail`, `search`, `uploads`, `shares`, `profile`, `settings`, `usage`, `applications`, `application-detail`, **`pair-login` (camera + envelope)** — plus the auth screens' submit paths. The **snake_case config** (`build.yaml`) must land before any model is built (README's "the one that silently breaks").
- **`26-flows` depends on the screens it chains.**
- **`27-platform-config` depends on `01-setup`; the biometric task in `04-core-infra` depends on its `FlutterFragmentActivity` note being real.**

## External edges (from `architecture.md` §7)

- **Mobile owns its own `flutter create`** — `flutter create --org com.bloberry --project-name bloberry mobile`. The one genuinely separate build unit (a Flutter project can't live inside a Go module's build).
- **Depends on `backend/tasks/19-openapi.md`** for `mobile/lib/api/` (the generated Dart client, PRD D8) — the client is generated from `api/openapi.yaml`; mobile's own `01-setup` doesn't hand-write it.
- **Flutter Desktop is NOT in scope for Bloberry** (desktop is Wails wrapping web, per TRD) — so no `flutter create --platforms=macos,windows,linux` addition; `desktop/tasks/` does not touch this project.

## Gaps flagged

None. All 21 screens have real mockups with Interactions; the style guide is complete; `mobile/navigation.md` has the route table + flow chains; the README pins packages, permissions and the queue design.
