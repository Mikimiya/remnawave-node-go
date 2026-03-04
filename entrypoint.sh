#!/bin/sh
set -e

GEO_DIR="/usr/local/share/xray"
GEOIP="$GEO_DIR/geoip.dat"
GEOSITE="$GEO_DIR/geosite.dat"

XRAY_VERSION="${XRAY_VERSION:-latest}"

detect_arch() {
    case "$(uname -m)" in
        'i386'|'i686')      echo '32' ;;
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
        *)
            echo "[Entrypoint] ERROR: unsupported architecture: $(uname -m)" >&2
            exit 1 ;;
    esac
}

mkdir -p "$GEO_DIR"

if [ ! -f "$GEOIP" ] || [ ! -f "$GEOSITE" ]; then
    echo "[Entrypoint] Geodata not found — downloading (xray-core ${XRAY_VERSION})..."

    ARCH="$(detect_arch)"

    if [ "$XRAY_VERSION" = "latest" ]; then
        DOWNLOAD_URL="https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-${ARCH}.zip"
    else
        DOWNLOAD_URL="https://github.com/XTLS/Xray-core/releases/download/${XRAY_VERSION}/Xray-linux-${ARCH}.zip"
    fi

    TMP_DIR="$(mktemp -d)"
    ZIP_FILE="${TMP_DIR}/xray.zip"

    echo "[Entrypoint] Fetching ${DOWNLOAD_URL}..."
    curl -sSfRL -H 'Cache-Control: no-cache' "$DOWNLOAD_URL" -o "$ZIP_FILE" || {
        echo "[Entrypoint] ERROR: download failed" >&2
        rm -rf "$TMP_DIR"
        exit 1
    }

    echo "[Entrypoint] Extracting geodata..."
    unzip -qj "$ZIP_FILE" "geoip.dat" "geosite.dat" -d "$TMP_DIR" || {
        echo "[Entrypoint] ERROR: extraction failed" >&2
        rm -rf "$TMP_DIR"
        exit 1
    }

    install -m 644 "$TMP_DIR/geoip.dat"   "$GEOIP"
    install -m 644 "$TMP_DIR/geosite.dat" "$GEOSITE"
    rm -rf "$TMP_DIR"

    echo "[Entrypoint] Geodata ready."
else
    echo "[Entrypoint] Geodata already present (xray-core ${XRAY_VERSION}), skipping download."
fi

exec "$@"
