## 1. Version Embedding

- [x] 1.1 Add version variable in main.go with ldflags placeholder
- [x] 1.2 Add `--version` flag handling to display version
- [x] 1.3 Update Makefile to inject version during build

## 2. GoReleaser Setup

- [x] 2.1 Install GoReleaser (`brew install goreleaser` or download)
- [x] 2.2 Create `.goreleaser.yaml` config with 5 platform targets
- [x] 2.3 Configure binary naming as `tsk_<os>_<arch>`
- [x] 2.4 Configure checksum generation (SHA256)
- [x] 2.5 Test local build with `goreleaser release --snapshot --clean`

## 3. GitHub Actions Workflow

- [x] 3.1 Create `.github/workflows/release.yml`
- [x] 3.2 Configure trigger on `v*` tags
- [x] 3.3 Add GoReleaser action step
- [x] 3.4 Configure GitHub token for release creation
- [x] 3.5 Add Homebrew tap update step

## 4. Homebrew Tap

- [x] 4.1 Configure main repo as Homebrew tap (Formula/ directory)
- [x] 4.2 Create initial `Formula/` directory 
- [x] 4.3 Configure GoReleaser to auto-update formula on release
- [x] 4.4 Test `brew tap Coleim/tsk && brew install tsk`

## 5. Install Script

- [x] 5.1 Create `install.sh` with platform detection (uname)
- [x] 5.2 Add GitHub API call to fetch latest release
- [x] 5.3 Implement download with curl/wget fallback
- [x] 5.4 Add SHA256 checksum verification
- [x] 5.5 Implement install location logic (/usr/local/bin vs ~/.local/bin)
- [x] 5.6 Add version argument support (`-s -- v1.2.3`)
- [x] 5.7 Add progress output messages
- [x] 5.8 Test on macOS

## 6. Testing & Documentation

- [x] 6.1 Create first release tag to test full pipeline
- [x] 6.2 Verify all platform binaries in GitHub Release
- [ ] 6.3 Test `go install github.com/coliva/tsk@latest`
- [x] 6.4 Test Homebrew installation flow
- [x] 6.5 Test curl installer
- [x] 6.6 Update README with installation instructions
