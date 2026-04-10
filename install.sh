#!/bin/sh
# tsk installer script
# Usage: curl -sSL https://raw.githubusercontent.com/Coleim/tsk/main/install.sh | bash
#    or: curl -sSL ... | bash -s -- v1.2.3   (specific version)

set -e

REPO="Coleim/tsk"
BINARY_NAME="tsk"
GITHUB_API="https://api.github.com/repos/${REPO}/releases"

# Colors (only if terminal supports it)
if [ -t 1 ]; then
    RED='\033[0;31m'
    GREEN='\033[0;32m'
    YELLOW='\033[0;33m'
    BLUE='\033[0;34m'
    NC='\033[0m' # No Color
else
    RED=''
    GREEN=''
    YELLOW=''
    BLUE=''
    NC=''
fi

info() {
    printf "${BLUE}==>${NC} %s\n" "$1"
}

success() {
    printf "${GREEN}==>${NC} %s\n" "$1"
}

warn() {
    printf "${YELLOW}Warning:${NC} %s\n" "$1"
}

error() {
    printf "${RED}Error:${NC} %s\n" "$1" >&2
    exit 1
}

# Detect OS
detect_os() {
    case "$(uname -s)" in
        Darwin)
            echo "darwin"
            ;;
        Linux)
            echo "linux"
            ;;
        MINGW*|MSYS*|CYGWIN*)
            echo "windows"
            ;;
        *)
            error "Unsupported operating system: $(uname -s)
Supported: macOS (darwin), Linux
For Windows, download directly from: https://github.com/${REPO}/releases"
            ;;
    esac
}

# Detect architecture
detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)
            echo "amd64"
            ;;
        arm64|aarch64)
            echo "arm64"
            ;;
        *)
            error "Unsupported architecture: $(uname -m)
Supported: x86_64 (amd64), arm64 (aarch64)"
            ;;
    esac
}

# Get download command (curl or wget)
get_downloader() {
    if command -v curl >/dev/null 2>&1; then
        echo "curl"
    elif command -v wget >/dev/null 2>&1; then
        echo "wget"
    else
        error "Neither curl nor wget found. Please install one of them."
    fi
}

# Download file
download() {
    url="$1"
    output="$2"
    downloader=$(get_downloader)
    
    case "$downloader" in
        curl)
            curl -fsSL "$url" -o "$output"
            ;;
        wget)
            wget -q "$url" -O "$output"
            ;;
    esac
}

# Fetch content (for API calls)
fetch() {
    url="$1"
    downloader=$(get_downloader)
    
    case "$downloader" in
        curl)
            curl -fsSL "$url"
            ;;
        wget)
            wget -qO- "$url"
            ;;
    esac
}

# Get latest version from GitHub API
get_latest_version() {
    version=$(fetch "${GITHUB_API}/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
    if [ -z "$version" ]; then
        error "Failed to fetch latest version from GitHub.
This might be due to API rate limiting. Try again later or specify a version:
  curl -sSL ... | bash -s -- v1.0.0"
    fi
    echo "$version"
}

# Verify checksum
verify_checksum() {
    archive="$1"
    expected_checksum="$2"
    
    if command -v sha256sum >/dev/null 2>&1; then
        actual_checksum=$(sha256sum "$archive" | awk '{print $1}')
    elif command -v shasum >/dev/null 2>&1; then
        actual_checksum=$(shasum -a 256 "$archive" | awk '{print $1}')
    else
        warn "Neither sha256sum nor shasum found. Skipping checksum verification."
        return 0
    fi
    
    if [ "$actual_checksum" != "$expected_checksum" ]; then
        error "Checksum verification failed!
Expected: $expected_checksum
Actual:   $actual_checksum
The downloaded file may be corrupted. Please try again."
    fi
}

# Get install directory
get_install_dir() {
    # Try /usr/local/bin first (requires sudo on most systems)
    if [ -w "/usr/local/bin" ]; then
        echo "/usr/local/bin"
        return
    fi
    
    # Fall back to ~/.local/bin
    local_bin="$HOME/.local/bin"
    mkdir -p "$local_bin"
    echo "$local_bin"
}

# Check if directory is in PATH
check_path() {
    dir="$1"
    case ":$PATH:" in
        *":$dir:"*)
            return 0
            ;;
        *)
            return 1
            ;;
    esac
}

main() {
    # Parse version argument
    VERSION="${1:-}"
    
    info "Detecting platform..."
    OS=$(detect_os)
    ARCH=$(detect_arch)
    info "Platform: ${OS}/${ARCH}"
    
    # Get version
    if [ -z "$VERSION" ]; then
        info "Fetching latest version..."
        VERSION=$(get_latest_version)
    fi
    info "Version: ${VERSION}"
    
    # Construct download URL
    ARCHIVE_NAME="${BINARY_NAME}_${OS}_${ARCH}.tar.gz"
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${ARCHIVE_NAME}"
    CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${VERSION}/checksums.txt"
    
    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap 'rm -rf "$TMP_DIR"' EXIT
    
    # Download archive
    info "Downloading ${ARCHIVE_NAME}..."
    download "$DOWNLOAD_URL" "$TMP_DIR/$ARCHIVE_NAME" || error "Failed to download from $DOWNLOAD_URL"
    
    # Download and verify checksum
    info "Verifying checksum..."
    download "$CHECKSUMS_URL" "$TMP_DIR/checksums.txt" || warn "Could not download checksums file"
    
    if [ -f "$TMP_DIR/checksums.txt" ]; then
        EXPECTED_CHECKSUM=$(grep "$ARCHIVE_NAME" "$TMP_DIR/checksums.txt" | awk '{print $1}')
        if [ -n "$EXPECTED_CHECKSUM" ]; then
            verify_checksum "$TMP_DIR/$ARCHIVE_NAME" "$EXPECTED_CHECKSUM"
            success "Checksum verified"
        else
            warn "Checksum for $ARCHIVE_NAME not found in checksums file"
        fi
    fi
    
    # Extract
    info "Extracting..."
    tar -xzf "$TMP_DIR/$ARCHIVE_NAME" -C "$TMP_DIR"
    
    # Install
    INSTALL_DIR=$(get_install_dir)
    info "Installing to ${INSTALL_DIR}..."
    
    if [ "$INSTALL_DIR" = "/usr/local/bin" ] && [ ! -w "$INSTALL_DIR" ]; then
        sudo mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        mv "$TMP_DIR/$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Verify installation
    if [ -x "$INSTALL_DIR/$BINARY_NAME" ]; then
        success "Successfully installed ${BINARY_NAME} ${VERSION} to ${INSTALL_DIR}/${BINARY_NAME}"
        
        # Check PATH
        if ! check_path "$INSTALL_DIR"; then
            warn "${INSTALL_DIR} is not in your PATH"
            echo ""
            echo "Add it to your shell profile:"
            echo "  export PATH=\"\$PATH:${INSTALL_DIR}\""
            echo ""
        fi
        
        # Show version
        echo ""
        "$INSTALL_DIR/$BINARY_NAME" --version
    else
        error "Installation failed"
    fi
}

main "$@"
