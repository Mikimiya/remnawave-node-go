#!/bin/bash
set -e

# ============================================================
#  Remnawave Node Go — Install Script (Binary)
#  https://github.com/hteppl/remnawave-node-go
#
#  Install:
#    bash <(curl -fsSL https://raw.githubusercontent.com/hteppl/remnawave-node-go/master/install.sh)
#
#  Update:
#    bash <(curl -fsSL https://raw.githubusercontent.com/hteppl/remnawave-node-go/master/install.sh) update
#
#  Update geodata only:
#    bash <(curl -fsSL https://raw.githubusercontent.com/hteppl/remnawave-node-go/master/install.sh) update-geo
#
#  Uninstall:
#    bash <(curl -fsSL https://raw.githubusercontent.com/hteppl/remnawave-node-go/master/install.sh) uninstall
# ============================================================

# --- Colors --------------------------------------------------
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

# --- Constants -----------------------------------------------
REPO_OWNER="hteppl"
REPO_NAME="remnawave-node-go"
GITHUB_REPO="${REPO_OWNER}/${REPO_NAME}"
GITHUB_RAW="https://raw.githubusercontent.com/${GITHUB_REPO}/master"
GITHUB_API="https://api.github.com/repos/${GITHUB_REPO}"

BINARY_NAME="remnawave-node-go"
INSTALL_BIN="/usr/local/bin/${BINARY_NAME}"
CONFIG_DIR="/etc/${BINARY_NAME}"
ENV_FILE="${CONFIG_DIR}/.env"
GEO_DIR="/usr/local/share/xray"
SERVICE_NAME="${BINARY_NAME}"
SERVICE_FILE="/etc/systemd/system/${SERVICE_NAME}.service"

# --- Helpers -------------------------------------------------
info()    { echo -e "${BLUE}[INFO]${NC}    $*"; }
success() { echo -e "${GREEN}[  OK]${NC}    $*"; }
warn()    { echo -e "${YELLOW}[WARN]${NC}    $*"; }
error()   { echo -e "${RED}[ERROR]${NC}   $*" >&2; }

banner() {
    echo -e "${CYAN}${BOLD}"
    cat << 'EOF'
  ____                                             _   _           _
 |  _ \ ___ _ __ ___  _ __   __ ___      _____   _| \ | | ___   __| | ___
 | |_) / _ \ '_ ` _ \| '_ \ / _` \ \ /\ / / _` |  \| |/ _ \ / _` |/ _ \
 |  _ <  __/ | | | | | | | | (_| |\ V  V / (_| | |\  | (_) | (_| |  __/
 |_| \_\___|_| |_| |_|_| |_|\__,_| \_/\_/ \__,_|_| \_|\___/ \__,_|\___|
                                              Go Edition — Binary Installer
EOF
    echo -e "${NC}"
}

check_root() {
    if [ "$(id -u)" -ne 0 ]; then
        error "This script must be run as root. Try:"
        error "  sudo bash install.sh"
        error "  sudo bash <(curl -fsSL ${GITHUB_RAW}/install.sh)"
        exit 1
    fi
}

# --- Detect platform -----------------------------------------
detect_platform() {
    OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
    case "$OS" in
        linux) ;;
        *)     error "Unsupported OS: $OS (only Linux is supported)"; exit 1 ;;
    esac

    ARCH="$(uname -m)"
    case "$ARCH" in
        x86_64|amd64)    PLATFORM_SUFFIX="linux-amd64" ;;
        aarch64|arm64)   PLATFORM_SUFFIX="linux-arm64" ;;
        armv7l|armv7)    PLATFORM_SUFFIX="linux-armv7"  ;;
        *)               error "Unsupported architecture: $ARCH"; exit 1 ;;
    esac

    info "Platform: ${PLATFORM_SUFFIX}"
}

# --- Detect xray arch for geodata ----------------------------
detect_xray_arch() {
    case "$(uname -m)" in
        'i386'|'i686')     echo '32' ;;
        'amd64'|'x86_64')  echo '64' ;;
        'armv5tel')         echo 'arm32-v5' ;;
        'armv6l')
            grep -qw 'vfp' /proc/cpuinfo 2>/dev/null && echo 'arm32-v6' || echo 'arm32-v5' ;;
        'armv7'|'armv7l')
            grep -qw 'vfp' /proc/cpuinfo 2>/dev/null && echo 'arm32-v7a' || echo 'arm32-v5' ;;
        'armv8'|'aarch64') echo 'arm64-v8a' ;;
        'mips')             echo 'mips32' ;;
        'mipsle')           echo 'mips32le' ;;
        'mips64')
            lscpu 2>/dev/null | grep -q 'Little Endian' && echo 'mips64le' || echo 'mips64' ;;
        'mips64le')         echo 'mips64le' ;;
        'ppc64')            echo 'ppc64' ;;
        'ppc64le')          echo 'ppc64le' ;;
        'riscv64')          echo 'riscv64' ;;
        's390x')            echo 's390x' ;;
        *)  error "Unsupported architecture for geodata: $(uname -m)"; exit 1 ;;
    esac
}

# --- Get latest release tag ----------------------------------
get_latest_version() {
    local tag
    tag=$(curl -sSf "${GITHUB_API}/releases/latest" 2>/dev/null \
        | grep '"tag_name"' | head -1 \
        | sed -E 's/.*"tag_name":\s*"([^"]+)".*/\1/')
    if [ -z "$tag" ]; then
        error "Failed to get latest version from ${GITHUB_REPO}."
        error "Check your network or visit: https://github.com/${GITHUB_REPO}/releases"
        exit 1
    fi
    echo "$tag"
}

# --- Install system dependencies -----------------------------
install_deps() {
    info "Checking dependencies..."
    local need=()

    for cmd in curl unzip; do
        command -v "$cmd" &>/dev/null || need+=("$cmd")
    done

    if [ ${#need[@]} -gt 0 ]; then
        info "Installing: ${need[*]}"
        if command -v apt-get &>/dev/null; then
            apt-get update -qq && apt-get install -y -qq "${need[@]}"
        elif command -v dnf &>/dev/null; then
            dnf install -y "${need[@]}"
        elif command -v yum &>/dev/null; then
            yum install -y "${need[@]}"
        elif command -v apk &>/dev/null; then
            apk add --no-cache "${need[@]}"
        elif command -v pacman &>/dev/null; then
            pacman -Sy --noconfirm "${need[@]}"
        else
            error "Cannot auto-install dependencies. Please install manually: ${need[*]}"
            exit 1
        fi
    fi

    success "Dependencies OK"
}

# --- Download prebuilt binary from GitHub Releases ------------
download_binary() {
    local version="$1"
    local binary_file="${BINARY_NAME}-${PLATFORM_SUFFIX}"
    local download_url="https://github.com/${GITHUB_REPO}/releases/download/${version}/${binary_file}"

    info "Downloading ${binary_file} (${version})..."

    local tmp
    tmp=$(mktemp)
    curl -sSfL "${download_url}" -o "$tmp" || {
        error "Download failed: ${download_url}"
        error "Make sure the release exists at: https://github.com/${GITHUB_REPO}/releases/tag/${version}"
        rm -f "$tmp"
        exit 1
    }

    # Verify sha256 checksum if available
    local sha_url="${download_url}.sha256"
    local sha_tmp
    sha_tmp=$(mktemp)
    if curl -sSfL "${sha_url}" -o "$sha_tmp" 2>/dev/null; then
        local expected actual
        expected=$(awk '{print $1}' "$sha_tmp")
        actual=$(sha256sum "$tmp" | awk '{print $1}')
        if [ "$expected" != "$actual" ]; then
            error "Checksum verification failed!"
            error "  Expected: ${expected}"
            error "  Actual:   ${actual}"
            rm -f "$tmp" "$sha_tmp"
            exit 1
        fi
        success "Checksum verified ✔"
    fi
    rm -f "$sha_tmp"

    install -m 755 "$tmp" "${INSTALL_BIN}"
    rm -f "$tmp"

    success "Binary installed → ${INSTALL_BIN} (${version})"
}

# --- Geodata -------------------------------------------------
download_geodata() {
    local force="${1:-false}"
    mkdir -p "$GEO_DIR"

    if [ "$force" = "true" ] || [ ! -f "${GEO_DIR}/geoip.dat" ] || [ ! -f "${GEO_DIR}/geosite.dat" ]; then
        info "Downloading geodata files..."
        local xarch
        xarch=$(detect_xray_arch)
        local url="https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-${xarch}.zip"

        local tmp
        tmp=$(mktemp -d)
        curl -sSfRL "$url" -o "${tmp}/xray.zip" || { error "Geodata download failed"; rm -rf "$tmp"; exit 1; }
        unzip -qjo "${tmp}/xray.zip" "geoip.dat" "geosite.dat" -d "$tmp" || { error "Geodata extract failed"; rm -rf "$tmp"; exit 1; }

        install -m 644 "${tmp}/geoip.dat"  "${GEO_DIR}/geoip.dat"
        install -m 644 "${tmp}/geosite.dat" "${GEO_DIR}/geosite.dat"
        rm -rf "$tmp"

        success "Geodata updated → ${GEO_DIR}"
    else
        success "Geodata already present (use 'update-geo' to force refresh)"
    fi
}

# --- Systemd service -----------------------------------------
create_service() {
    info "Creating systemd service..."

    cat > "${SERVICE_FILE}" << EOF
[Unit]
Description=Remnawave Node Go
Documentation=https://github.com/${GITHUB_REPO}
After=network.target network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=${INSTALL_BIN}
EnvironmentFile=${ENV_FILE}
Environment=XRAY_LOCATION_ASSET=${GEO_DIR}
WorkingDirectory=${CONFIG_DIR}
Restart=always
RestartSec=5
LimitNOFILE=65535
StandardOutput=journal
StandardError=journal
SyslogIdentifier=${BINARY_NAME}

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    systemctl enable "${SERVICE_NAME}" &>/dev/null

    success "Service created and enabled"
}

# --- Interactive config --------------------------------------
configure_env() {
    mkdir -p "${CONFIG_DIR}"

    if [ -f "${ENV_FILE}" ]; then
        warn "Config already exists → ${ENV_FILE} (keeping)"
        return
    fi

    echo ""
    echo -e "${BOLD}──────────────────────── Configuration ────────────────────────${NC}"
    echo ""

    if [ -n "$SECRET_KEY" ]; then
        info "Using SECRET_KEY from environment"
    else
        while true; do
            read -rp "$(echo -e "${CYAN}SECRET_KEY (Base64 from Remnawave panel): ${NC}")" SECRET_KEY
            [ -n "$SECRET_KEY" ] && break
            warn "SECRET_KEY is required."
        done
    fi

    read -rp "$(echo -e "${CYAN}NODE_PORT [2222]: ${NC}")" NODE_PORT
    NODE_PORT="${NODE_PORT:-2222}"

    read -rp "$(echo -e "${CYAN}INTERNAL_REST_PORT [61001]: ${NC}")" INTERNAL_REST_PORT
    INTERNAL_REST_PORT="${INTERNAL_REST_PORT:-61001}"

    read -rp "$(echo -e "${CYAN}LOG_LEVEL (debug/info/warn/error) [info]: ${NC}")" LOG_LEVEL
    LOG_LEVEL="${LOG_LEVEL:-info}"

    cat > "${ENV_FILE}" << EOF
# Remnawave Node Go — generated $(date -u '+%Y-%m-%d %H:%M:%S UTC')
SECRET_KEY=${SECRET_KEY}
NODE_PORT=${NODE_PORT}
INTERNAL_REST_PORT=${INTERNAL_REST_PORT}
LOG_LEVEL=${LOG_LEVEL}
EOF
    chmod 600 "${ENV_FILE}"

    success "Config saved → ${ENV_FILE}"
}

# --- Print summary -------------------------------------------
print_done() {
    local action="${1:-installed}"
    echo ""
    echo -e "${BOLD}═══════════════════════════════════════════════════════════════${NC}"
    echo -e "  ${GREEN}${BOLD}✔ Remnawave Node Go ${action} successfully!${NC}"
    echo -e "${BOLD}═══════════════════════════════════════════════════════════════${NC}"
    echo ""
    echo -e "  ${GREEN}Binary  :${NC}  ${INSTALL_BIN}"
    echo -e "  ${GREEN}Config  :${NC}  ${ENV_FILE}"
    echo -e "  ${GREEN}Geodata :${NC}  ${GEO_DIR}"
    echo -e "  ${GREEN}Service :${NC}  ${SERVICE_NAME}"
    echo ""
    echo -e "  ${CYAN}Service management:${NC}"
    echo -e "    systemctl start   ${SERVICE_NAME}"
    echo -e "    systemctl stop    ${SERVICE_NAME}"
    echo -e "    systemctl restart ${SERVICE_NAME}"
    echo -e "    systemctl status  ${SERVICE_NAME}"
    echo ""
    echo -e "  ${CYAN}View logs:${NC}"
    echo -e "    journalctl -u ${SERVICE_NAME} -f"
    echo ""
    echo -e "  ${CYAN}Edit config:${NC}"
    echo -e "    nano ${ENV_FILE}"
    echo ""
}

# =============================================================
#  Commands
# =============================================================

cmd_install() {
    check_root
    detect_platform
    install_deps

    local version
    version=$(get_latest_version)
    info "Latest release: ${version}"
    echo ""

    download_binary "$version"
    download_geodata
    configure_env
    create_service

    info "Starting service..."
    systemctl restart "${SERVICE_NAME}"
    sleep 2

    if systemctl is-active --quiet "${SERVICE_NAME}"; then
        success "Service is running"
    else
        warn "Service may not have started. Check: journalctl -u ${SERVICE_NAME} -e"
    fi

    print_done "installed"
}

cmd_update() {
    check_root
    detect_platform
    install_deps

    local version
    version=$(get_latest_version)
    info "Updating to: ${version}"
    echo ""

    if systemctl is-active --quiet "${SERVICE_NAME}" 2>/dev/null; then
        info "Stopping service..."
        systemctl stop "${SERVICE_NAME}"
    fi

    download_binary "$version"
    download_geodata
    create_service

    info "Starting service..."
    systemctl restart "${SERVICE_NAME}"
    sleep 2

    if systemctl is-active --quiet "${SERVICE_NAME}"; then
        success "Service is running"
    else
        warn "Service may not have started. Check: journalctl -u ${SERVICE_NAME} -e"
    fi

    print_done "updated"
}

cmd_update_geo() {
    check_root
    detect_platform

    download_geodata "true"

    if systemctl is-active --quiet "${SERVICE_NAME}" 2>/dev/null; then
        info "Restarting service..."
        systemctl restart "${SERVICE_NAME}"
        success "Service restarted with updated geodata"
    fi
    echo ""
}

cmd_uninstall() {
    check_root

    echo ""
    warn "This will remove the Remnawave Node Go binary, service, and geodata."
    read -rp "$(echo -e "${RED}${BOLD}Proceed with uninstall? (y/N): ${NC}")" CONFIRM
    [[ ! "$CONFIRM" =~ ^[Yy]$ ]] && { info "Aborted."; exit 0; }

    # Stop & remove service
    if systemctl is-active --quiet "${SERVICE_NAME}" 2>/dev/null; then
        info "Stopping service..."
        systemctl stop "${SERVICE_NAME}"
    fi
    if [ -f "${SERVICE_FILE}" ]; then
        systemctl disable "${SERVICE_NAME}" &>/dev/null || true
        rm -f "${SERVICE_FILE}"
        systemctl daemon-reload
        success "Service removed"
    fi

    # Remove binary
    [ -f "${INSTALL_BIN}" ] && { rm -f "${INSTALL_BIN}"; success "Binary removed"; }

    # Remove geodata
    [ -d "${GEO_DIR}" ] && { rm -rf "${GEO_DIR}"; success "Geodata removed"; }

    # Config
    if [ -d "${CONFIG_DIR}" ]; then
        read -rp "$(echo -e "${YELLOW}Also remove config (${CONFIG_DIR})? (y/N): ${NC}")" RM_CONF
        if [[ "$RM_CONF" =~ ^[Yy]$ ]]; then
            rm -rf "${CONFIG_DIR}"
            success "Config removed"
        else
            info "Config preserved → ${CONFIG_DIR}"
        fi
    fi

    echo ""
    success "Remnawave Node Go uninstalled."
    echo ""
}

# =============================================================
#  Entry
# =============================================================
main() {
    banner

    case "${1:-install}" in
        install)     cmd_install    ;;
        update)      cmd_update     ;;
        update-geo)  cmd_update_geo ;;
        uninstall)   cmd_uninstall  ;;
        *)
            echo "Usage: bash install.sh {install|update|update-geo|uninstall}"
            echo ""
            echo "  install      Install Remnawave Node Go (default)"
            echo "  update       Update binary to latest version"
            echo "  update-geo   Force-update geodata files only"
            echo "  uninstall    Stop service and remove everything"
            echo ""
            exit 1
            ;;
    esac
}

main "$@"
