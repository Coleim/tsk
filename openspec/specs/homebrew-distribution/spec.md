## ADDED Requirements

### Requirement: Main repo serves as Homebrew tap
The main repository SHALL contain a `Formula/` directory for Homebrew formulas.

#### Scenario: Tap can be added
- **WHEN** user runs `brew tap coliva/tsk`
- **THEN** the tap is added successfully

### Requirement: Formula installs tsk
The formula SHALL download and install the correct binary for the user's platform.

#### Scenario: Install on macOS arm64
- **WHEN** user runs `brew install coliva/tsk/tsk` on macOS arm64
- **THEN** the darwin/arm64 binary is downloaded and installed

#### Scenario: Install on macOS amd64
- **WHEN** user runs `brew install coliva/tsk/tsk` on macOS amd64
- **THEN** the darwin/amd64 binary is downloaded and installed

#### Scenario: Install on Linux
- **WHEN** user runs `brew install coliva/tsk/tsk` on Linux
- **THEN** the appropriate linux binary is downloaded and installed

### Requirement: Formula validates checksum
The formula SHALL verify the SHA256 checksum of downloaded binaries.

#### Scenario: Checksum verification
- **WHEN** binary is downloaded during install
- **THEN** SHA256 checksum is verified against expected value

### Requirement: Upgrade support
The formula SHALL support version upgrades via standard Homebrew commands.

#### Scenario: Upgrade to new version
- **WHEN** new version is released and user runs `brew upgrade tsk`
- **THEN** the new version is downloaded and installed
