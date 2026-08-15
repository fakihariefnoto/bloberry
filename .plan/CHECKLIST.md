# Checklist — Bloberry

## Design

- [x] PRD written (`PRD.md`)
- [x] PRD detailed (`detail-prd`) — specific Problem/Goals, explicit Non-goals, prioritized user stories, Open questions resolved/sharpened
- [x] Architecture written (`architecture.md`, `detail-architecture`) — context + container diagrams, cross-platform flow sequences, deployment topology, cross-cutting concerns, ADR-style decision records. Runs after the PRD, before the TRD.
- [x] Implementation layout decided (`architecture.md` §7) — the concrete repo tree, the plan-folder→real-path mapping, and **which platform's setup task owns the single project `init`**. Every `build-*` skill reads this before writing its setup file; without it an embedded backend or a desktop shell gets scaffolded once per plan folder.
- [x] TRD written (`TRD.md`)
- [x] TRD detailed (`detail-trd`) — real architecture description, concrete data model, API/integration list, real technical risks
- [x] Style guide complete (`design/style-guide.md`) — every scale filled in (color light+dark, type with line-heights, spacing, radius, elevation, sizing/touch targets, motion) and component specs include disabled/pressed/focus/error states. No "TBD" — a gap here becomes an inconsistency on every screen.
- [x] Navigation map written (`mobile/navigation.md` / `web/navigation.md`) — route graph + route table (name, path, typed args, presentation, back, auth), flow chains with cancel/failure landings
- [x] Mobile mockups (`mobile/mockup/`) — purpose/states/real content per screen, not placeholder boxes
- [x] Web mockups (`web/mockup/`, mobile + desktop) — purpose/states/real content per screen, not placeholder boxes
- [x] Mockups generated for every screen (`generate-mockups`), not just the first one
- [x] Navigation graph closed — no exit points to a screen that doesn't exist, no entry unmirrored by a real exit, no unreachable screen, every mockup in the route table
- [x] CLI commands designed (`cli/commands/`, `generate-commands`) — real `--help` text and sample output per state, not just a command list. Only if this app has CLI.
- [x] Backend plan (`backend/`)
- [x] CLI plan (`cli/`) — only if this app has CLI
- [x] Desktop plan (`desktop/`) — only if this app has Desktop
- [x] Infra plan (`infra/`)

## Mobile (Flutter) — only if this app has Mobile

Defaults defined in `templates/flutter-defaults.md`; see `mobile/README.md`.

- [x] App identifier set (`TRD.md`) — `com.bloberry.app` and identical on Android `applicationId` and iOS bundle identifier
- [x] Required screens planned: Onboarding (one-time), Welcome, Signup, Login, Profile, Settings (or explicit reason noted for any that are skipped)
- [x] UI kit picked (shadcn/ui for Flutter, Forui, Material, or GetWidget) and recorded in `design/style-guide.md`
- [ ] Loading state + disabled-while-processing guard applied to every heavy/async action
- [ ] No overflow-pixel warnings on any screen
- [ ] HTTP responses parsed via the standard envelope (`{data?, messages?}`), no raw package/exception text shown to users
- [x] JSON deserialization configured for snake_case (`fieldRename: FieldRename.snake` or equivalent) — not left at Dart's camelCase default
- [x] Platform permissions listed per native capability used (AndroidManifest.xml + Info.plist)
- [x] PIN/biometric unlock available as an optional Settings toggle
- [x] If biometric: `MainActivity` extends `FlutterFragmentActivity` (Android), and its auth failure path is logged, not silently swallowed
- [ ] Native/platform-channel failure paths use `debugPrint`, not silent catches
- [x] Form validators (email/password/phone) centralized in one file, not duplicated per screen
- [x] Phone number input (if any) — N/A, no phone input in this app uses a maintained dial-code+flag package, not a hand-rolled list

## Web — only if this app has Web

Defaults defined in `templates/web-defaults.md`; see `web/README.md`.

- [x] JS runtime / package manager confirmed — Bun (Node + pnpm default, or Bun / Deno)
- [x] Dev/staging/prod env files set up, `.env.example` committed, build scripts wired per environment
- [x] (Vue only) Nuxt.js vs plain Vue decided, and UI kit picked (PrimeVue/shadcn-vue/Ant Design Vue/Element Plus/Vuetify/Reka UI), recorded in `design/style-guide.md`
- [x] Shared/reusable components identified across planned screens (`web/components.md`), not copy-pasted per page

## Backend — only if this app has Backend

Defaults defined in `templates/backend-{go,python,java}-defaults.md` (matching the chosen language); see `backend/README.md`.

- [x] Web framework picked, recorded in `TRD.md` — chi
- [x] Database(s) picked, recorded in `TRD.md` (real database — never an in-memory mock)
- [x] User roles need decided (yes/no); if yes, role list + schema/collection design recorded
- [x] Domain structure planned (`backend/domains.md`) — layered (handler/usecase/repository or language equivalent) with interface injection between layers
- [x] Auth domain covers signup/login/refresh/logout, forgot password, login by OTP, and login with Google
- [x] User domain covers profile + settings, per the standard pattern
- [x] Session/refresh-token store decided (Redis by default) with platform-aware TTLs (mobile vs. web) — ephemeral data only; durable data lives in the chosen database
- [x] OpenAPI spec workflow confirmed, committed spec file path set, cross-referenced from `mobile/README.md`/`web/README.md`
- [x] Interface naming consistent per the language convention (mock-friendly)
- [x] `make lint` and `make security` guard targets present in `backend/Makefile`

## CLI — only if this app has CLI

Defaults defined in `templates/cli-defaults.md` and `templates/cli-distribution.md`; see `cli/README.md`.

- [x] Language picked (Go / Rust / Python / Node+TS), recorded in `TRD.md` — Go
- [x] Role decided (standalone, or companion to the Backend — sharing its auth and OpenAPI contract), recorded in `TRD.md`
- [x] Command tree written (`cli/README.md`) and confirmed with the user before individual commands were designed
- [x] Global flags settled (`--help`, `--version`, `--config`, `--json`, `--quiet`, `--verbose`, `--no-color`, `--yes`)
- [x] Exit-code table written — scripts branch on these, so they're API, not an implementation detail
- [x] Config file location + precedence order (flag → env → file → default) recorded
- [x] Secrets stored in the OS keychain or an env var — never plaintext in the config file
- [x] (Companion) Login shape picked (browser device/auth-code flow vs personal access token) and `auth logout` revokes server-side
- [x] stdout-is-data / stderr-is-everything-else confirmed, and `--json` coverage decided per command
- [x] Color, spinners and prompts gated on TTY; every prompt has a flag equivalent so the tool works in CI
- [x] Shell completions planned (bash/zsh/fish/powershell) **and** installed by the packages, not just generated
- [x] Distribution table written (`cli/README.md`) — one row per channel with its exact install command, host, automation and credentials
- [x] Release automation picked (GoReleaser) (GoReleaser / cargo-dist / fpm) so per-release manifest bumps aren't manual
- [ ] Install docs written — the README section listing each channel's command
- [x] Update story decided, including that self-update never overwrites a Homebrew/apt-managed binary

## Desktop — only if this app has Desktop

Defaults defined in `templates/desktop-defaults.md`; see `desktop/README.md`.

- [x] Framework picked (Tauri / Electron / Wails / Flutter Desktop), recorded in `TRD.md` — Wails
- [x] Wrapped/reused frontend confirmed (Web's build output, or Mobile's Flutter codebase) and that platform detailed first
- [x] Framework sub-choices settled (frontend source + IPC + auto-updater for Tauri/Wails; contextIsolation/nodeIntegration baseline + auto-updater for Electron; window-management package for Flutter Desktop)
- [x] App menu structure, keyboard shortcuts, and window-state persistence recorded
- [x] System tray in scope decided explicitly (yes — background uploads/sync need it) (yes with reason, or no)
- [x] Mockup plan confirmed as reuse-not-redraw (chrome-only deltas in `desktop/README.md`, not a parallel screen set)
- [x] Packaging/signing/notarization plan set per target OS (or explicitly deferred, noted as such)
- [x] CI build-matrix exception flagged in `infra/README.md` (macOS/Windows builds need OS-native CI, not the self-hosted VPS runner)

## Infra — only if this app has Infra

Defaults defined in `templates/infra-defaults.md`; see `infra/README.md`. GitHub Actions with a self-hosted runner and manual-only triggers (including CI checks) is a fixed convention — not re-decided per app.

- [x] VPS process model picked per component — systemd (Docker Compose / systemd), recorded in `TRD.md`
- [x] Reverse proxy picked (Caddy / nginx) — Caddy; domains noted in `detail-infra`
- [x] `ci.yml` and `deploy.yml` generated from `templates/github-actions/` per component, secrets list matches each `.env.example`
- [x] GitHub Environments (`staging`, `production`) planned with matching secrets — noted as a manual step, not automatable from here
- [ ] Self-hosted runner registered and running as a service on the VPS

## Build

Each item is expanded into a concrete task list by its `build-*` skill — a numbered `<platform>/tasks/` folder (00-index.md + one file per group), see "Planning the build" in `AGENTS.md`. Check off here only once that platform's task folder is fully done, not when it's merely started.

- [ ] Mobile — see `mobile/tasks/` (`build-mobile`) **— task list generated (28 files), work not started**
- [ ] Backend — see `backend/tasks/` (`build-backend`) **— task list generated (21 files), work not started**
- [ ] Web — see `web/tasks/` (`build-web`) **— task list generated (36 files), work not started**
- [ ] CLI — only if this app has CLI; see `cli/tasks/` (`build-cli` — including the distribution channels, each verified by actually installing from it) **— task list generated (15 files), work not started**
- [ ] Desktop — only if this app has Desktop; see `desktop/tasks/` (`build-desktop` — chrome/window/packaging tasks only, since screen content is already covered by `mobile/tasks/`/`web/tasks/`) **— task list generated (9 files), work not started**
- [ ] Infra / deployment — see `infra/tasks/` (`build-infra`) **— task list generated (8 files), work not started**

## Launch

Expanded into a concrete task list by `build-launch` — see `LAUNCH.md`.

- [ ] Internal testing
- [ ] Release
