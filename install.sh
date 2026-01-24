#!/bin/bash

# Kitty-Beads Build and Install Script
# Builds from source and installs the binary
# Usage: curl -sSL https://raw.githubusercontent.com/AkimZmerli/kitty-beads/main/scripts/install-from-source.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
GITHUB_REPO="AkimZmerli/kitty-beads"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
TEMP_DIR=$(mktemp -d)

cleanup() {
    rm -rf "$TEMP_DIR"
}

trap cleanup EXIT

# Check prerequisites
check_prerequisites() {
    echo -e "${BLUE}Checking prerequisites...${NC}"

    local missing_deps=()

    if ! command -v git &> /dev/null; then
        missing_deps+=("git")
    fi

    if ! command -v go &> /dev/null; then
        missing_deps+=("go (1.24+)")
    else
        local go_version=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+' || echo "unknown")
        echo -e "${GREEN}✓ Go $go_version found${NC}"
    fi

    if [ ${#missing_deps[@]} -gt 0 ]; then
        echo -e "${RED}Error: Missing required dependencies:${NC}"
        for dep in "${missing_deps[@]}"; do
            echo "  - $dep"
        done
        echo ""
        echo -e "${YELLOW}Installation instructions:${NC}"
        echo "  macOS (homebrew): brew install git go"
        echo "  Ubuntu/Debian: sudo apt-get install git golang-go"
        echo "  Fedora: sudo dnf install git golang"
        echo ""
        echo "Then run this script again."
        exit 1
    fi

    echo -e "${GREEN}✓ All prerequisites met${NC}"
}

# Clone repository
clone_repo() {
    echo ""
    echo -e "${BLUE}Cloning repository...${NC}"

    git clone --depth 1 https://github.com/$GITHUB_REPO "$TEMP_DIR/kitty-beads"
    cd "$TEMP_DIR/kitty-beads"

    echo -e "${GREEN}✓ Repository cloned${NC}"
}

# Build binary
build_binary() {
    echo ""
    echo -e "${BLUE}Building kitty-beads...${NC}"

    if [ ! -f "Makefile" ]; then
        echo -e "${RED}Error: Makefile not found${NC}"
        exit 1
    fi

    make build

    if [ ! -f "bin/kitty-beads" ]; then
        echo -e "${RED}Error: Build failed - binary not created${NC}"
        exit 1
    fi

    echo -e "${GREEN}✓ Build successful${NC}"
}

# Install binary
install_binary() {
    echo ""
    echo -e "${BLUE}Installing to: $INSTALL_DIR/kitty-beads${NC}"

    mkdir -p "$INSTALL_DIR"

    if [ -w "$INSTALL_DIR" ]; then
        cp "$TEMP_DIR/kitty-beads/bin/kitty-beads" "$INSTALL_DIR/kitty-beads"
        chmod +x "$INSTALL_DIR/kitty-beads"
        echo -e "${GREEN}✓ Installed to: $INSTALL_DIR/kitty-beads${NC}"
    else
        echo -e "${YELLOW}Installation directory requires sudo${NC}"
        sudo cp "$TEMP_DIR/kitty-beads/bin/kitty-beads" "$INSTALL_DIR/kitty-beads"
        sudo chmod +x "$INSTALL_DIR/kitty-beads"
        echo -e "${GREEN}✓ Installed to: $INSTALL_DIR/kitty-beads${NC}"
    fi
}

# Verify installation
verify_installation() {
    echo ""
    echo -e "${BLUE}Verifying installation...${NC}"

    if command -v kitty-beads &> /dev/null; then
        echo -e "${GREEN}✓ kitty-beads is available in PATH${NC}"
        kitty-beads --version 2>/dev/null || echo -e "${GREEN}✓ Binary verified${NC}"
    else
        echo -e "${YELLOW}⚠ kitty-beads not in PATH, but installed to: $INSTALL_DIR/kitty-beads${NC}"
        echo "Add to your shell config:"
        echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    fi
}

# Show next steps
next_steps() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}✓ Kitty-Beads built and installed successfully!${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "${BLUE}Quick start:${NC}"
    echo ""
    echo "1. Go to your project:"
    echo -e "   ${YELLOW}cd your-project${NC}"
    echo ""
    echo "2. Initialize Beads (if not done):"
    echo -e "   ${YELLOW}bd init${NC}"
    echo ""
    echo "3. Start kitty-beads dashboard:"
    echo -e "   ${YELLOW}kitty-beads${NC}"
    echo ""
    echo "4. Open http://localhost:8080"
    echo ""
    echo "5. Create your first issue:"
    echo -e "   ${YELLOW}bd create \"Your first task\"${NC}"
    echo ""
    echo -e "${BLUE}Available commands:${NC}"
    echo "  kitty-beads              - Start dashboard on :8080"
    echo "  kitty-beads -port 3000   - Custom port"
    echo "  kitty-beads --help       - Show options"
    echo ""
    echo -e "${BLUE}Useful resources:${NC}"
    echo "  GitHub: https://github.com/$GITHUB_REPO"
    echo "  Beads: https://github.com/steveyegge/beads"
    echo ""
}

# Main flow
main() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║          Kitty-Beads Build & Install Script               ║${NC}"
    echo -e "${BLUE}║         Building from latest source code                  ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""

    check_prerequisites
    clone_repo
    build_binary
    install_binary
    verify_installation
    next_steps
}

main
