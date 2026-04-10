## ADDED Requirements

### Requirement: One-line curl installation
Users SHALL be able to install via `curl -sSL <url> | bash`.

#### Scenario: Basic installation
- **WHEN** user runs `curl -sSL https://raw.githubusercontent.com/coliva/tsk/main/install.sh | bash`
- **THEN** the latest tsk binary is downloaded and installed

### Requirement: Platform auto-detection
The script SHALL detect the user's OS and architecture automatically.

#### Scenario: Detect macOS arm64
- **WHEN** script runs on macOS with Apple Silicon
- **THEN** darwin/arm64 binary is selected

#### Scenario: Detect Linux amd64
- **WHEN** script runs on Linux x86_64
- **THEN** linux/amd64 binary is selected

#### Scenario: Unsupported platform error
- **WHEN** script runs on unsupported platform (e.g., FreeBSD)
- **THEN** clear error message shows supported platforms and manual download URL

### Requirement: Install location with fallback
The script SHALL install to `/usr/local/bin` if writable, otherwise `~/.local/bin`.

#### Scenario: System-wide install with sudo
- **WHEN** user has write access to `/usr/local/bin`
- **THEN** binary is installed to `/usr/local/bin/tsk`

#### Scenario: User-local install without sudo
- **WHEN** user cannot write to `/usr/local/bin`
- **THEN** binary is installed to `~/.local/bin/tsk` and PATH warning shown if needed

### Requirement: Version specification
Users SHALL be able to install a specific version.

#### Scenario: Install specific version
- **WHEN** user runs `curl -sSL <url> | bash -s -- v1.2.3`
- **THEN** version v1.2.3 is installed instead of latest

### Requirement: Checksum verification
The script SHALL verify the SHA256 checksum of downloaded binary.

#### Scenario: Checksum passes
- **WHEN** binary is downloaded
- **THEN** SHA256 is verified against checksums.txt from the release

#### Scenario: Checksum fails
- **WHEN** downloaded binary checksum doesn't match
- **THEN** installation aborts with error message

### Requirement: Clear progress output
The script SHALL show clear progress messages during installation.

#### Scenario: Installation progress
- **WHEN** script runs
- **THEN** output shows: detecting platform, downloading, verifying, installing, success message with version
