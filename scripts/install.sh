#!/usr/bin/env bash
#
# gitloom install script
# Installs gitloom binary for Linux, macOS, and Windows
#

set -euo pipefail

# Configuration
GITHUB_REPO="ovitorvalente/git-loom"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"
INSTALL_MODE="${INSTALL_MODE:-auto}"
FORCE="${FORCE:-false}"
UNINSTALL="${UNINSTALL:-false}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_banner() {
    echo ""
    echo -e "${BOLD}  ██████╗ ███████╗███████╗██╗     ██╗███╗   ██╗███████╗";
    echo -e "  ██╔══██╗██╔════╝██╔════╝██║     ██║████╗  ██║██╔════╝";
    echo -e "  ██║  ██║█████╗  █████╗  ██║     ██║██╔██╗ ██║█████╗  ";
    echo -e "  ██║  ██║██╔══╝  ██╔══╝  ██║     ██║██║╚██╗██║██╔══╝  ";
    echo -e "  ██████╔╝███████╗██║     ███████╗██║██║ ╚████║███████╗";
    echo -e "  ╚═════╝ ╚══════╝╚═╝     ╚══════╝╚═╝╚═╝  ╚═══╝╚══════╝${NC}";
    echo -e "${BOLD}  Smart Git Commits${NC}"
    echo ""
}

usage() {
    cat << EOF
Usage: install.sh [OPTIONS]

Install gitloom binary

OPTIONS:
    -v, --version VERSION    Install specific version (default: latest)
    -d, --dir DIR            Install directory (default: /usr/local/bin)
    -m, --mode MODE          Install mode: auto, user, system (default: auto)
    -f, --force              Force reinstall
    -u, --uninstall          Uninstall gitloom
    -h, --help               Show this help

EXAMPLES:
    # Install latest version
    curl -fsSL https://raw.githubusercontent.com/${GITHUB_REPO}/main/scripts/install.sh | bash

    # Install specific version
    curl -fsSL https://raw.githubusercontent.com/${GITHUB_REPO}/main/scripts/install.sh | bash -s -- -v 0.1.0-alpha

    # Install to custom directory
    curl -fsSL https://raw.githubusercontent.com/${GITHUB_REPO}/main/scripts/install.sh | bash -s -- -d ~/.local/bin

    # Uninstall
    curl -fsSL https://raw.githubusercontent.com/${GITHUB_REPO}/main/scripts/install.sh | bash -s -- -u

MODES:
    auto    - Install system-wide if root, otherwise user directory
    user    - Install to ~/.local/bin
    system  - Install to /usr/local/bin (requires root)
EOF
    exit 0
}

get_os() {
    local os
    os="$(uname -s | tr '[:upper:]' '[:lower:]')"
    case "$os" in
        linux*)  echo "linux" ;;
        darwin*) echo "darwin" ;;
        msys*|mingw*|cygwin*) echo "windows" ;;
        *) log_error "Unsupported OS: $os"; exit 1 ;;
    esac
}

get_arch() {
    local arch
    arch="$(uname -m | tr '[:upper:]' '[:lower:]')"
    case "$arch" in
        x86_64|amd64)   echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        armv7l|armhf)   echo "armv7" ;;
        *) log_error "Unsupported architecture: $arch"; exit 1 ;;
    esac
}

get_latest_version() {
    local version
    version=$(curl -sSL "https://api.github.com/repos/${GITHUB_REPO}/releases/latest" | grep -o '"tag_name":.*' | sed 's/.*": "//' | sed 's/",//' | sed 's/^v//')
    echo "${version}"
}

get_version_for_filename() {
    local tag_version="$1"
    echo "$tag_version" | sed 's/-.*$//'
}

check_curl() {
    if ! command -v curl &> /dev/null; then
        log_error "curl is required but not installed"
        exit 1
    fi
}

check_deps() {
    if command -v sha256sum &> /dev/null; then
        echo "sha256sum"
    elif command -v shasum &> /dev/null; then
        echo "shasum -a 256"
    elif command -v sha256 &> /dev/null; then
        echo "sha256"
    else
        log_warn "No checksum tool found, skipping verification"
        echo ""
    fi
}

download_checksum() {
    local version="$1"
    local checksum_url="https://github.com/${GITHUB_REPO}/releases/download/v${version}/checksums.txt"
    curl -sSL "$checksum_url" 2>/dev/null || echo ""
}

verify_checksum() {
    local file="$1"
    local expected_checksum="$2"
    
    if [ -z "$expected_checksum" ]; then
        log_warn "No checksum provided, skipping verification"
        return 0
    fi
    
    local actual_checksum
    actual_checksum=$($(check_deps) "$file" 2>/dev/null | awk '{print $1}')
    
    if [ "$actual_checksum" != "$expected_checksum" ]; then
        log_error "Checksum mismatch!"
        log_error "Expected: $expected_checksum"
        log_error "Actual:   $actual_checksum"
        rm -f "$file"
        exit 1
    fi
    
    log_success "Checksum verified"
}

get_download_url() {
    local version="$1"
    local version_filename="$2"
    local os="$3"
    local arch="$4"
    
    local ext="tar.gz"
    if [ "$os" = "windows" ]; then
        ext="zip"
    fi
    
    echo "https://github.com/${GITHUB_REPO}/releases/download/v${version}/gitloom_v${version_filename}_${os}_${arch}.${ext}"
}

do_install() {
    local version="${1:-}"
    local force="${2:-false}"
    
    local os
    local arch
    local install_dir
    local bin_name="gitloom"
    local version_for_filename
    
    os=$(get_os)
    arch=$(get_arch)
    
    if [ -z "$version" ]; then
        log_info "Fetching latest version..."
        version=$(get_latest_version)
    fi
    
    version_for_filename=$(get_version_for_filename "$version")
    
    log_info "Version: ${BOLD}${version}${NC}"
    log_info "OS: ${BOLD}${os}${NC}"
    log_info "Arch: ${BOLD}${arch}${NC}"
    
    # Determine install directory
    case "$INSTALL_MODE" in
        auto)
            if [ "$(id -u)" -eq 0 ]; then
                install_dir="$INSTALL_DIR"
            else
                install_dir="${HOME}/.local/bin"
            fi
            ;;
        user)
            install_dir="${HOME}/.local/bin"
            ;;
        system)
            install_dir="$INSTALL_DIR"
            ;;
        *)
            install_dir="$INSTALL_DIR"
            ;;
    esac
    
    log_info "Install dir: ${BOLD}${install_dir}${NC}"
    
    # Check if already installed
    if [ -f "${install_dir}/${bin_name}" ] && [ "$force" != "true" ]; then
        local current_version
        current_version=$("${install_dir}/${bin_name}" version 2>/dev/null | head -1 | awk '{print $2}' || echo "unknown")
        
        if [ "$current_version" = "$version" ]; then
            log_info "gitloom ${version} is already installed"
            exit 0
        fi
        
        log_warn "gitloom ${current_version} is installed. Upgrading to ${version}..."
    fi
    
    # Create temp directory
    local tmp_dir
    tmp_dir=$(mktemp -d)
    trap "rm -rf $tmp_dir" EXIT
    
    # Download
    local download_url
    local checksum_url
    local archive_file
    local ext="tar.gz"
    
    if [ "$os" = "windows" ]; then
        ext="zip"
    fi
    
    download_url=$(get_download_url "$version" "$version_for_filename" "$os" "$arch")
    archive_file="${tmp_dir}/gitloom.${ext}"
    
    log_info "Downloading..."
    log_info "  ${download_url}"
    
    if ! curl -fSL "$download_url" -o "$archive_file"; then
        log_error "Failed to download gitloom"
        log_error "URL: $download_url"
        log_error "Version may not exist for this platform"
        exit 1
    fi
    
    # Download checksum
    log_info "Downloading checksums..."
    local checksums
    checksums=$(download_checksum "$version" "$os" "$arch")
    
    # Verify
    if [ -n "$checksums" ]; then
        local expected_sum
        expected_sum=$(echo "$checksums" | grep "${os}_${arch}\.${ext}" | awk '{print $1}' || echo "")
        verify_checksum "$archive_file" "$expected_sum"
    fi
    
    # Extract
    log_info "Installing..."
    mkdir -p "$install_dir"
    
    local archive_dir="gitloom_v${version_for_filename}_${os}_${arch}"
    
    if [ "$ext" = "zip" ]; then
        unzip -o "$archive_file" -d "$tmp_dir" > /dev/null
        mv "${tmp_dir}/${archive_dir}/gitloom.exe" "${install_dir}/${bin_name}.exe" 2>/dev/null || \
        unzip -j "$archive_file" "*/gitloom.exe" -d "$install_dir" > /dev/null
        chmod +x "${install_dir}/${bin_name}.exe"
    else
        tar xzf "$archive_file" -C "$tmp_dir"
        mv "${tmp_dir}/${archive_dir}/gitloom" "${install_dir}/${bin_name}"
        chmod +x "${install_dir}/${bin_name}"
    fi
    
    # Verify installation
    if [ -f "${install_dir}/${bin_name}" ] || [ -f "${install_dir}/${bin_name}.exe" ]; then
        local final_bin="${install_dir}/${bin_name}"
        [ "$os" = "windows" ] && final_bin="${install_dir}/${bin_name}.exe"
        
        log_success "Installed ${BOLD}${final_bin}${NC}"
    else
        log_error "Installation failed"
        exit 1
    fi
    
    # Check PATH
    if [[ ":$PATH:" != *":${install_dir}:"* ]]; then
        log_warn "${install_dir} is not in your PATH"
        
        if [ "$INSTALL_MODE" = "user" ] || [ "$INSTALL_MODE" = "auto" ]; then
            echo ""
            echo "Add to your shell profile (~/.bashrc, ~/.zshrc, etc.):"
            echo -e "  ${GREEN}export PATH=\"\$HOME/.local/bin:\$PATH\"${NC}"
        fi
    fi
    
    echo ""
    log_success "Done! Run ${BOLD}gitloom --help${NC} to get started"
}

do_uninstall() {
    local bin_name="gitloom"
    local removed=false
    
    # Try common locations
    local locations=(
        "/usr/local/bin/${bin_name}"
        "/usr/bin/${bin_name}"
        "${HOME}/.local/bin/${bin_name}"
        "${HOME}/.local/bin/${bin_name}.exe"
    )
    
    for loc in "${locations[@]}"; do
        if [ -f "$loc" ]; then
            log_info "Removing $loc"
            rm -f "$loc"
            removed=true
        fi
    done
    
    if [ "$removed" = true ]; then
        log_success "gitloom uninstalled"
    else
        log_warn "gitloom not found in standard locations"
    fi
}

main() {
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case "$1" in
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            -d|--dir)
                INSTALL_DIR="$2"
                shift 2
                ;;
            -m|--mode)
                INSTALL_MODE="$2"
                shift 2
                ;;
            -f|--force)
                FORCE="true"
                shift
                ;;
            -u|--uninstall)
                UNINSTALL="true"
                shift
                ;;
            -h|--help)
                usage
                ;;
            *)
                log_error "Unknown option: $1"
                usage
                ;;
        esac
    done
    
    print_banner
    
    check_curl
    
    if [ "$UNINSTALL" = "true" ]; then
        do_uninstall
    else
        do_install "${VERSION:-}" "${FORCE:-false}"
    fi
}

main "$@"
