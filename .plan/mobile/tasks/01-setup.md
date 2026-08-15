# Task group — 01 setup

**Depends on:** nothing (own project). **Blocks:** everything in this folder.

**Context from `architecture.md` §7:** mobile is the one genuinely separate build unit — its own `flutter create`. Desktop is Wails wrapping web, so Flutter Desktop platforms are NOT added.

- [ ] **`flutter create --org com.bloberry --project-name bloberry mobile`** at the repo root — the project's one init.
- [ ] **App identifier matches `TRD.md` exactly** — `com.bloberry.app` as Android `applicationId` (`android/app/build.gradle`) and iOS `PRODUCT_BUNDLE_IDENTIFIER` (`ios/Runner.xcodeproj`). The exact string on both, per the README.
- [ ] **Packages added** from `mobile/README.md` §Packages: `dio`, `json_serializable` + `build_runner`, `flutter_secure_storage`, `local_auth`, `image_picker`, `file_picker`, `path_provider` + `sqflite`, `connectivity_plus`, `flutter_spinkit`, `awesome_snackbar_content`, `share_plus`, `cached_network_image`, `go_router`, `fl_chart`. **Deliberately not**: maps, calendar_view, timelines.
- [ ] **`shadcn/ui for Flutter`** added per the style guide (`design/style-guide.md` UI-kit line) — the composable kit driven by the shared tokens.
- [ ] **Folder tree** per §7: `mobile/lib/{screens,widgets,api,store,theme}/` — screens (one per route), widgets (shared), api (the generated Dart client), store (upload queue + session), theme (the token output of `02`).
- [ ] **`Makefile`** — build/run/debug/test targets from `templates/makefiles/flutter.mk`, adjusted for this project layout.
- [ ] **`build.yaml`** — `json_serializable` configured with `field_rename: snake` **globally** (README's snake_case trap — not per-class annotations). This is verified before any model is built.
- [ ] **Smoke run** — `flutter run` boots to a placeholder shell, `flutter test` passes.

**tests:** `flutter analyze` clean; `flutter test` green on the scaffold.
