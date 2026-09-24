#!/bin/sh
# Собирает один .deb пакет wb-mqtt-noolite из уже собранных бинарников.
# Бинарник статический (CGO_ENABLED=0) - для Debian 12 (bookworm) и 13 (trixie) содержимое
# пакета идентично, различается только версия (суффикс ~bookworm/~trixie) для того, чтобы
# оба варианта можно было держать в apt-репозитории одновременно, в отдельных suite.
#
# Использование:
#   packaging/build-deb.sh <debian-arch> <suite> <path-to-binaries> <version> <output-dir>
# Пример:
#   packaging/build-deb.sh armhf bookworm build/armv7 2.0.0 build/deb

set -eu

if [ "$#" -ne 5 ]; then
    echo "Usage: $0 <debian-arch> <suite> <bin-dir> <version> <out-dir>" >&2
    exit 1
fi

DEB_ARCH="$1"
SUITE="$2"
BIN_DIR="$3"
VERSION="$4"
OUT_DIR="$5"

PKG_NAME="wb-mqtt-noolite"
PKG_VERSION="${VERSION}~${SUITE}"
MAINTAINER="${MAINTAINER:-Sergey Butenin <s.butenin@gmail.com>}"

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
PKG_DIR="${ROOT_DIR}/packaging/deb"

for tool in dpkg-deb fakeroot; do
    if ! command -v "$tool" >/dev/null 2>&1; then
        echo "Требуется $tool (пакет dpkg-dev/fakeroot)" >&2
        exit 1
    fi
done

for bin in wb-mqtt-noolite mtrf_tool; do
    if [ ! -f "${BIN_DIR}/${bin}" ]; then
        echo "Не найден ${BIN_DIR}/${bin} - сначала выполните 'make armv7' или 'make arm64'" >&2
        exit 1
    fi
done

STAGE="$(mktemp -d)"
trap 'rm -rf "$STAGE"' EXIT
chmod 0755 "$STAGE"

mkdir -p "$STAGE/DEBIAN" \
         "$STAGE/usr/bin" \
         "$STAGE/etc" \
         "$STAGE/lib/systemd/system" \
         "$STAGE/usr/share/doc/$PKG_NAME"

install -m 0755 "$BIN_DIR/wb-mqtt-noolite" "$STAGE/usr/bin/wb-mqtt-noolite"
install -m 0755 "$BIN_DIR/mtrf_tool" "$STAGE/usr/bin/mtrf_tool"

install -m 0644 "$PKG_DIR/wb-mqtt-noolite.json" "$STAGE/etc/wb-mqtt-noolite.json"
install -m 0644 "$ROOT_DIR/templates/templates.json" "$STAGE/etc/wb-mqtt-noolite-templates.json"
install -m 0644 "$PKG_DIR/wb-mqtt-noolite-devices.json" "$STAGE/etc/wb-mqtt-noolite-devices.json"

install -m 0644 "$ROOT_DIR/packaging/wb-mqtt-noolite.service" "$STAGE/lib/systemd/system/wb-mqtt-noolite.service"
install -m 0644 "$PKG_DIR/copyright" "$STAGE/usr/share/doc/$PKG_NAME/copyright"

install -m 0755 "$PKG_DIR/postinst" "$STAGE/DEBIAN/postinst"
install -m 0755 "$PKG_DIR/prerm" "$STAGE/DEBIAN/prerm"
install -m 0755 "$PKG_DIR/postrm" "$STAGE/DEBIAN/postrm"

cat > "$STAGE/DEBIAN/conffiles" <<EOF
/etc/wb-mqtt-noolite.json
/etc/wb-mqtt-noolite-templates.json
/etc/wb-mqtt-noolite-devices.json
EOF

SIZE_KB=$(du -sk "$STAGE" | cut -f1)

cat > "$STAGE/DEBIAN/control" <<EOF
Package: $PKG_NAME
Version: $PKG_VERSION
Section: net
Priority: optional
Architecture: $DEB_ARCH
Installed-Size: $SIZE_KB
Maintainer: $MAINTAINER
Recommends: mosquitto
Homepage: https://github.com/basicus/wb-mqtt-noolite
Description: Noolite(-F) to MQTT bridge for Wiren Board
 wb-mqtt-noolite integrates Noolite(-F) devices (remotes, switches and
 sensors) into a Wiren Board automation controller over the MTRF-64-USB-A
 adapter, following the Wiren Board MQTT Conventions and WB-STD-001.
 .
 Includes mtrf_tool, a command-line utility for pairing/unpairing devices
 with the adapter.
EOF

mkdir -p "$OUT_DIR"
OUT_FILE="${OUT_DIR}/${PKG_NAME}_${PKG_VERSION}_${DEB_ARCH}.deb"
fakeroot dpkg-deb --root-owner-group --build "$STAGE" "$OUT_FILE" >&2

echo "$OUT_FILE"
