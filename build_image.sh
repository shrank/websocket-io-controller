#!/bin/bash
set -e

NAME=io2websocket-gateway
VERSION=${VERSION:-0.0.1}
BUILD=${BUILD:-$(date +%Y%m%d%H%M)}
PKG=${NAME}-${VERSION}-${BUILD}

TC_IMAGE=$1
TC_OUT=${TC_OUT:-tc-${NAME}-${VERSION}-${BUILD}.img}
OVERLAY_OUT=${OVERLAY_OUT:-overlay-${NAME}-${VERSION}-${BUILD}.gz}

GOOS=${GOOS:-linux}
GOARCH=${GOARCH:-arm64}
GOARM=${GOARM:-}

CWD=$(pwd)
WORK_DIR=${CWD}/build/${PKG}
OVERLAY_DIR=${WORK_DIR}/overlay
ROOT_DIR=${OVERLAY_DIR}/root

log() {
    echo "[build] $*"
}

cleanup() {
    if [ -d "${WORK_DIR}" ]; then
        rm -rf "${WORK_DIR}"
    fi
}

cleanup
mkdir -p "${ROOT_DIR}/usr/local/bin"
mkdir -p "${ROOT_DIR}/usr/local/etc/init.d"
mkdir -p "${ROOT_DIR}/usr/local/share/${NAME}/public"
mkdir -p "${ROOT_DIR}/opt"

log "building frontend ui"
cd frontend
npm install
npm run build
cd "${CWD}"

log "cross-compiling ${NAME} for raspberry pi (${GOOS}/${GOARCH}${GOARM:+v${GOARM}})"
GOOS=${GOOS} GOARCH=${GOARCH} ${GOARM:+GOARM=${GOARM}} go build -o "${ROOT_DIR}/usr/local/bin/${NAME}" .

log "installing config and frontend assets"
cp config.txt "${ROOT_DIR}/usr/local/etc/${NAME}.conf"
cp -r frontend/dist/. "${ROOT_DIR}/usr/local/share/${NAME}/public/"

cat > "${ROOT_DIR}/usr/local/etc/init.d/${NAME}" <<EOF
#!/bin/sh
# Tiny Core boot script for io2websocket-gateway

case "\$1" in
    start)
        if [ -f /usr/local/etc/${NAME}.conf ]; then
            CONFIG=/usr/local/etc/${NAME}.conf
        else
            CONFIG=/opt/${NAME}.conf
        fi
        echo "starting ${NAME}..."
        /usr/local/bin/${NAME} -debug-dir /usr/local/share/${NAME}/public "\$CONFIG" > /var/log/${NAME}.log 2>&1 &
        ;;
    stop)
        echo "stopping ${NAME}..."
        killall ${NAME} 2>/dev/null || true
        ;;
    restart)
        \$0 stop
        sleep 1
        \$0 start
        ;;
    *)
        echo "usage: \$0 {start|stop|restart}"
        exit 1
        ;;
esac
EOF
chmod +x "${ROOT_DIR}/usr/local/etc/init.d/${NAME}"

log "generating bootlocal.sh hook"
mkdir -p "${ROOT_DIR}/opt"
cat > "${ROOT_DIR}/opt/bootlocal.sh" <<EOF
#!/bin/sh
/usr/local/etc/init.d/${NAME} start
EOF
chmod +x "${ROOT_DIR}/opt/bootlocal.sh"

log "packaging overlay as gzip compressed cpio"
cd "${ROOT_DIR}"
find . -print | cpio -o -H newc | gzip -9 > "${CWD}/build/${OVERLAY_OUT}"
cd "${CWD}"

log "overlay created: build/${OVERLAY_OUT}"

if [ -n "${TC_IMAGE}" ] && [ -f "${TC_IMAGE}" ]; then
    log "embedding overlay into ${TC_IMAGE} -> ${TC_OUT}"
    cp "${TC_IMAGE}" "${CWD}/build/${TC_OUT}"

    LOOP_DEV=$(sudo losetup -fP --show "${CWD}/build/${TC_OUT}")
    MOUNT_POINT=$(mktemp -d)

    sudo mount "${LOOP_DEV}p1" "${MOUNT_POINT}" || sudo mount "${LOOP_DEV}" "${MOUNT_POINT}"
    if [ ! -d "${MOUNT_POINT}/overlays" ]; then
        sudo umount "${MOUNT_POINT}"
        sudo losetup -d "${LOOP_DEV}"
        rmdir "${MOUNT_POINT}"
        log "error: ${MOUNT_POINT}/overlays not found, is this a Tiny Core image?"
        exit 1
    fi

    sudo cp "${CWD}/build/${OVERLAY_OUT}" "${MOUNT_POINT}/${OVERLAY_OUT}"

    sudo sed -i -e 's/\.gz /.gz,${OVERLAY_OUT} /g' ${MOUNT_POINT}/config.txt
    sudo sh -c "echo dtparam=i2c_arm=on >> ${MOUNT_POINT}/config.txt"
    sudo sh -c "echo dtparam=spi=on >> ${MOUNT_POINT}/config.txt"

    sudo umount "${MOUNT_POINT}"
    sudo losetup -d "${LOOP_DEV}"
    rmdir "${MOUNT_POINT}"
    rm -r ${WORK_DIR}
    log "embedded image created: build/${TC_OUT}"
else
    log "TC_IMAGE not provided or not found, skipping image embedding"
fi

log "build complete"
