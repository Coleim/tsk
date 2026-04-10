## Why

Currently `tsk` can only be installed by cloning the repo and running `go build`. Users need Go installed and must manage updates manually. Making the app easily installable via common methods (`go install`, `brew`, `curl`) increases adoption and simplifies updates.

## What Changes

- Add GitHub Actions release pipeline that builds cross-platform binaries
- Create Homebrew tap formula for macOS/Linux installation
- Create install.sh script for quick curl-based installation
- Configure module for `go install` support
- Set up GitHub Releases as distribution point

## Capabilities

### New Capabilities
- `release-pipeline`: GitHub Actions workflow for automated releases, cross-platform builds, and artifact publishing
- `homebrew-distribution`: Homebrew tap formula and tap repository setup
- `script-installer`: curl-installable shell script with version detection and platform support

### Modified Capabilities
<!-- No existing spec-level behavior changes needed -->

## Impact

- **New files**: `.github/workflows/release.yml`, `install.sh`, Homebrew formula
- **External**: Requires GitHub repository setup for releases, separate tap repo for Homebrew
- **Dependencies**: GoReleaser or similar for cross-compilation
- **Platforms**: macOS (arm64, amd64), Linux (arm64, amd64), Windows (amd64)
