## ADDED Requirements

### Requirement: Release triggered by version tag
The release workflow SHALL be triggered when a tag matching `v*` pattern is pushed.

#### Scenario: Tag push triggers release
- **WHEN** a tag `v1.2.3` is pushed to the repository
- **THEN** the release workflow starts automatically

#### Scenario: Non-version tags ignored
- **WHEN** a tag `latest` or `test-build` is pushed
- **THEN** no release workflow is triggered

### Requirement: Cross-platform binary builds
The workflow SHALL produce binaries for darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, and windows/amd64.

#### Scenario: All platforms built
- **WHEN** release workflow runs
- **THEN** binaries are produced for all 5 platform combinations

#### Scenario: Binary naming convention
- **WHEN** binaries are built
- **THEN** each binary is named `tsk_<os>_<arch>` (e.g., `tsk_darwin_arm64`)

### Requirement: Version embedded in binary
The binary SHALL report its version when run with `--version` flag.

#### Scenario: Version flag shows tag version
- **WHEN** user runs `tsk --version`
- **THEN** output shows the git tag version (e.g., `tsk version v1.2.3`)

### Requirement: GitHub Release created with artifacts
The workflow SHALL create a GitHub Release with all binaries attached.

#### Scenario: Release contains all binaries
- **WHEN** workflow completes successfully
- **THEN** a GitHub Release exists with all platform binaries attached

#### Scenario: Release includes checksums
- **WHEN** release is created
- **THEN** a `checksums.txt` file is attached with SHA256 hashes of all binaries

### Requirement: Homebrew formula updated on release
The workflow SHALL update the Homebrew tap formula after successful release.

#### Scenario: Formula auto-updated
- **WHEN** release workflow completes
- **THEN** the `homebrew-tap` repository formula is updated with new version and checksums
