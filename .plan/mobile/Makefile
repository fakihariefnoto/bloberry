# Mobile Makefile — Flutter
#
# Reference template copied into apps/bloberry/mobile/Makefile by the
# plan-app skill (mobile is always Flutter). Trimmed from a real project's
# Makefile — adjust BACKEND_PORT / staging URL / bundle IDs as needed.

# ─── Configuration ───────────────────────────────────────────
# Android emulator type: avd (Android Studio) -> 10.0.2.2, genymotion -> 10.0.3.2
ANDROID_EMULATOR ?= avd
ifeq ($(ANDROID_EMULATOR),avd)
  ANDROID_DEV_IP ?= 10.0.2.2
else ifeq ($(ANDROID_EMULATOR),genymotion)
  ANDROID_DEV_IP ?= 10.0.3.2
else
  ANDROID_DEV_IP ?= $(ANDROID_EMULATOR)
endif
BACKEND_PORT ?= 8000

APP_NAME     := $(shell grep -E '^name:' pubspec.yaml | awk '{print $$2}')
VERSION      ?= $(shell grep -E '^version:' pubspec.yaml | awk '{print $$2}')
BUILD_NAME   := $(word 1,$(subst +, ,$(VERSION)))
BUILD_NUMBER := $(word 2,$(subst +, ,$(VERSION)))
_build_flags = --build-name=$(BUILD_NAME) --build-number=$(BUILD_NUMBER)

_get_simulator = $(shell flutter devices 2>/dev/null | grep '(simulator)' | head -1 | awk -F'•' '{gsub(/[[:space:]]/, "", $$2); print $$2}')
_get_emulator  = $(shell flutter devices 2>/dev/null | grep -E '(emulator|mobile)' | head -1 | awk -F'•' '{gsub(/[[:space:]]/, "", $$2); print $$2}')

# ─── Dev / Debug ─────────────────────────────────────────────
.PHONY: run-ios
run-ios: ## Run on iOS simulator (calls localhost, hot reload + debugger attached)
	@open -a Simulator 2>/dev/null || true
	@sleep 3
	flutter run -d $(call _get_simulator) --dart-define=API_BASE_URL=http://localhost:$(BACKEND_PORT)

.PHONY: run-android
run-android: ## Run on Android emulator (calls $(ANDROID_DEV_IP), hot reload + debugger attached)
	flutter run -d $(call _get_emulator) --dart-define=API_BASE_URL=http://$(ANDROID_DEV_IP):$(BACKEND_PORT)

.PHONY: debug
debug: run-ios ## Alias for run-ios — `flutter run` attaches the Dart debugger by default

# ─── Build ────────────────────────────────────────────────────
.PHONY: build-android-debug
build-android-debug: ## Build Android APK (debug)
	flutter build apk --debug --dart-define=API_BASE_URL=http://$(ANDROID_DEV_IP):$(BACKEND_PORT)

.PHONY: build-ios-debug
build-ios-debug: ## Build iOS (debug, simulator)
	flutter build ios --debug --simulator --dart-define=API_BASE_URL=http://localhost:$(BACKEND_PORT)

.PHONY: build-android-release
build-android-release: ## Build Android AAB (release)
	flutter build appbundle --release $(_build_flags) --obfuscate --split-debug-info=build/debug-info

.PHONY: build-ios-release
build-ios-release: ## Build iOS (release)
	flutter build ios --release $(_build_flags)

# ─── Code generation ──────────────────────────────────────────
.PHONY: gen
gen: ## Run build_runner code generation
	dart run build_runner build --delete-conflicting-outputs

# ─── Testing / Quality ─────────────────────────────────────────
.PHONY: test
test: ## Run unit tests
	flutter test

.PHONY: test-coverage
test-coverage: ## Run unit tests with coverage
	flutter test --coverage

.PHONY: analyze
analyze: ## Run flutter analyze
	flutter analyze

.PHONY: format
format: ## Format all dart files
	dart format .

# ─── Dependencies / Clean ──────────────────────────────────────
.PHONY: get
get: ## Get dependencies
	flutter pub get

.PHONY: clean
clean: ## Clean build artifacts
	flutter clean

# ─── Help ────────────────────────────────────────────────────
.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
