## Context

`tsk` is a Go-based terminal task manager. Currently distributed as source code only. Module path is `github.com/coliva/tsk`. The app is a single binary with no external runtime dependencies, making it ideal for simple distribution.

## Goals / Non-Goals

**Goals:**
- Single command installation on macOS/Linux
- Automated release process triggered by git tags
- Cross-platform binaries for major OS/arch combinations
- Version embedded in binary

**Non-Goals:**
- Windows installer (.msi) - simple .exe download is sufficient
- Auto-update mechanism - users update via their install method
- Package manager support beyond Homebrew (apt, yum, etc.) - future consideration

## Decisions

### 1. Use GoReleaser for builds
**Decision**: Use GoReleaser for cross-compilation and release automation.
**Rationale**: Industry standard for Go projects, handles cross-compilation, checksums, changelogs. Single config file.
**Alternative considered**: Manual `go build` scripts - more maintenance, less features.

### 2. Homebrew Tap repository
**Decision**: Create `homebrew-tap` repository at `github.com/coliva/homebrew-tap`.
**Rationale**: Standard Homebrew tap naming convention. Can host multiple formulas if needed later. GoReleaser can auto-update the formula.
**Alternative considered**: Submit to homebrew-core - requires more maintenance, popularity threshold.

### 3. Install script fetches from GitHub Releases
**Decision**: Shell script downloads latest release from GitHub API, detects platform, installs to `/usr/local/bin` or `~/.local/bin`.
**Rationale**: Common pattern (rustup, homebrew). No external dependencies. Fallback to home directory if no sudo.
**Alternative considered**: Host on separate CDN - unnecessary complexity.

### 4. Semantic versioning with git tags
**Decision**: Release triggered by pushing `v*` tags (e.g., `v1.0.0`).
**Rationale**: Standard Go versioning. Works with `go install`, GoReleaser, and Homebrew.

## Risks / Trade-offs

**[Risk]** GitHub API rate limits for unauthenticated requests (60/hour)
→ Install script caches release info, provides manual download fallback

**[Risk]** Homebrew tap requires maintenance
→ GoReleaser auto-updates formula on release

**[Risk]** Platform detection in install script may fail on edge cases
→ Provide clear error messages with manual download instructions

**[Trade-off]** No Windows script installer
→ Windows users download .exe directly from releases - acceptable for v1
