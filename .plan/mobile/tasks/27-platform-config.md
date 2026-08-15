# Task group — 27 platform config

**Depends on:** `01-setup.md`. **Blocks:** nothing (finishing config).

- [ ] **Camera** — Android `android.permission.CAMERA`; iOS `NSCameraUsageDescription` — "Take photos to upload to your Bloberry folders, and scan QR codes to sign in." (The same permission serves capture and `pair-login`, M22 — one string covering both benefits.)
- [ ] **Photo library** — Android `READ_MEDIA_IMAGES`, `READ_MEDIA_VIDEO` (API 33+); iOS `NSPhotoLibraryUsageDescription` — "Select photos and videos to upload."
- [ ] **File selection** — none (scoped storage / document picker, per the README table).
- [ ] **Network** — Android `android.permission.INTERNET`.
- [ ] **Biometrics** — Android `android.permission.USE_BIOMETRIC`; iOS `NSFaceIDUsageDescription` — "Unlock Bloberry with Face ID."
- [ ] **Background upload** — Android `FOREGROUND_SERVICE`, `FOREGROUND_SERVICE_DATA_SYNC` (API 34+); iOS `UIBackgroundModes: [processing]`.
- [ ] **Purpose strings state the actual benefit** — generic strings get rejected by App Review and declined by users (README).
- [ ] **App icon + splash assets** — per platform, using the brand mark.
- [ ] **App name** — "Bloberry" on both launchers.

**tests:** `flutter build apk --debug` and the iOS equivalent compile with the manifest/plist changes; permission prompts show the real strings.
