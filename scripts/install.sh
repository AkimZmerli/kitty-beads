#!/bin/bash

# Kitty-Beads Installation Script
# Installs kitty-beads binary for your platform
# Usage: curl -sSL https://raw.githubusercontent.com/AkimZmerli/kitty-beads/main/scripts/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="AkimZmerli/kitty-beads"
INSTALL_DIR="${INSTALL_DIR:-.}"
BINARY_NAME="kitty-beads"

# Detect OS and architecture
detect_platform() {
    local os_type=$(uname -s)
    local arch=$(uname -m)

    case "$os_type" in
        Darwin)
            OS="darwin"
            ;;
        Linux)
            OS="linux"
            ;;
        *)
            echo -e "${RED}Error: Unsupported OS: $os_type${NC}"
            echo "Supported: macOS (Darwin), Linux"
            echo "For other platforms, build from source: git clone https://github.com/$GITHUB_REPO && make build"
            exit 1
            ;;
    esac

    case "$arch" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}Error: Unsupported architecture: $arch${NC}"
            echo "Supported: x86_64 (amd64), aarch64/arm64"
            exit 1
            ;;
    esac

    echo -e "${BLUE}Detected: $os_type ($arch)${NC}"
}

# Get latest release version
get_latest_version() {
    echo -e "${BLUE}Fetching latest release...${NC}"

    local latest_url="https://api.github.com/repos/$GITHUB_REPO/releases/latest"
    local version=$(curl -s "$latest_url" | grep '"tag_name"' | head -1 | cut -d'"' -f4)

    if [ -z "$version" ]; then
        echo -e "${YELLOW}Warning: Could not fetch latest release${NC}"
        echo -e "${YELLOW}Trying alternative: checking releases page...${NC}"

        # Fallback: try to get releases list
        local releases=$(curl -s "https://api.github.com/repos/$GITHUB_REPO/releases?per_page=1")
        version=$(echo "$releases" | grep '"tag_name"' | head -1 | cut -d'"' -f4)

        if [ -z "$version" ]; then
            echo -e "${RED}Error: Could not determine latest version${NC}"
            echo "Please check: https://github.com/$GITHUB_REPO/releases"
            exit 1
        fi
    fi

    RELEASE_VERSION="${version#v}"  # Remove 'v' prefix if present
    echo -e "${GREEN}Latest version: $RELEASE_VERSION${NC}"
}

# Download binary
download_binary() {
    local tmp_dir=$(mktemp -d)
    local archive_name="${BINARY_NAME}_${RELEASE_VERSION}_${OS}_${ARCH}.tar.gz"
    local download_url="https://github.com/$GITHUB_REPO/releases/download/v${RELEASE_VERSION}/${archive_name}"

    echo -e "${BLUE}Downloading from: $download_url${NC}"

    if ! curl -fsSL -o "$tmp_dir/$archive_name" "$download_url"; then
        # Try alternative formats
        local zip_name="${BINARY_NAME}_${RELEASE_VERSION}_${OS}_${ARCH}.zip"
        local zip_url="https://github.com/$GITHUB_REPO/releases/download/v${RELEASE_VERSION}/${zip_name}"

        echo -e "${YELLOW}TAR download failed, trying ZIP format...${NC}"
        if ! curl -fsSL -o "$tmp_dir/$zip_name" "$zip_url"; then
            echo -e "${RED}Error: Could not download binary${NC}"
            echo "Tried: $download_url"
            echo "Also tried: $zip_url"
            echo ""
            echo -e "${YELLOW}Alternative: Build from source${NC}"
            echo "  git clone https://github.com/$GITHUB_REPO"
            echo "  cd kitty-beads"
            echo "  make build"
            rm -rf "$tmp_dir"
            exit 1
        fi

        # Extract ZIP
        unzip -q "$tmp_dir/$zip_name" -d "$tmp_dir"
        BINARY_PATH="$tmp_dir/$BINARY_NAME"
    else
        # Extract TAR
        tar -xzf "$tmp_dir/$archive_name" -C "$tmp_dir"
        BINARY_PATH="$tmp_dir/$BINARY_NAME"
    fi

    if [ ! -f "$BINARY_PATH" ]; then
        echo -e "${RED}Error: Binary not found in archive${NC}"
        rm -rf "$tmp_dir"
        exit 1
    fi

    chmod +x "$BINARY_PATH"
    TMP_DIR="$tmp_dir"
    echo -e "${GREEN}Downloaded successfully${NC}"
}

# Install binary
install_binary() {
    echo -e "${BLUE}Installing to: $INSTALL_DIR/$BINARY_NAME${NC}"

    if [ "$INSTALL_DIR" = "." ]; then
        cp "$BINARY_PATH" "./$BINARY_NAME"
        echo -e "${GREEN}✓ Installed to: $(pwd)/$BINARY_NAME${NC}"
    else
        mkdir -p "$INSTALL_DIR"
        cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"

        # Try to add to PATH if installing to /usr/local/bin
        if [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
            echo -e "${GREEN}✓ Installed to: $INSTALL_DIR/$BINARY_NAME${NC}"

            # Verify it's in PATH
            if ! command -v kitty-beads &> /dev/null; then
                echo -e "${YELLOW}⚠ Warning: $INSTALL_DIR is not in your PATH${NC}"
                echo "Add to your shell config (~/.bashrc, ~/.zshrc, etc.):"
                echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
            fi
        else
            echo -e "${GREEN}✓ Installed to: $INSTALL_DIR/$BINARY_NAME${NC}"
            echo -e "${YELLOW}Make sure $INSTALL_DIR is in your PATH${NC}"
        fi
    fi
}

# Verify installation
verify_installation() {
    echo -e "${BLUE}Verifying installation...${NC}"

    local binary_to_check
    if [ "$INSTALL_DIR" = "." ]; then
        binary_to_check="./$BINARY_NAME"
    else
        binary_to_check="$INSTALL_DIR/$BINARY_NAME"
    fi

    if "$binary_to_check" --version &>/dev/null || "$binary_to_check" -version &>/dev/null; then
        echo -e "${GREEN}✓ Installation verified${NC}"
        "$binary_to_check" --version 2>/dev/null || true
    elif "$binary_to_check" --help &>/dev/null || "$binary_to_check" -help &>/dev/null; then
        echo -e "${GREEN}✓ Installation verified${NC}"
    else
        echo -e "${YELLOW}⚠ Warning: Could not verify binary${NC}"
        echo "But the file was installed successfully."
    fi
}

# Cleanup
cleanup() {
    if [ -n "$TMP_DIR" ] && [ -d "$TMP_DIR" ]; then
        rm -rf "$TMP_DIR"
    fi
}

# Show next steps
next_steps() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}✓ Kitty-Beads installed successfully!${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${BLUE}Next steps:${NC}"
    echo ""
    echo "1. Initialize a project (if not already done):"
    echo -e "   ${YELLOW}cd your-project${NC}"
    echo -e "   ${YELLOW}bd init${NC}"
    echo ""
    echo "2. Start the dashboard:"
    echo -e "   ${YELLOW}kitty-beads${NC}"
    echo -e "   (or: ${YELLOW}kitty-beads -port 3000${NC} for custom port)"
    echo ""
    echo "3. Open in your browser:"
    echo -e "   ${YELLOW}http://localhost:8080${NC}"
    echo ""
    echo "4. Create your first issue:"
    echo -e "   ${YELLOW}bd create \"Your first task\"${NC}"
    echo ""
    echo -e "${BLUE}Learn more:${NC}"
    echo "  Docs: https://github.com/$GITHUB_REPO"
    echo "  Beads: https://github.com/steveyegge/beads"
    echo ""
}

# Main installation flow
main() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║             Kitty-Beads Installer                          ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""

    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --install-dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            --help)
                echo "Usage: ./install.sh [OPTIONS]"
                echo ""
                echo "Options:"
                echo "  --install-dir DIR  Installation directory (default: current directory)"
                echo "  --help             Show this help message"
                exit 0
                ;;
            *)
                echo "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done

    detect_platform
    get_latest_version
    download_binary

    trap cleanup EXIT

    install_binary
    verify_installation
    next_steps
}

main "$@"
