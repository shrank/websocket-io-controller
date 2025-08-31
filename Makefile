PHONY: all frontend server dist io2websocket-gateway

VERSION="0.0.1"
BUILD=$(shell date +%Y%m%d%H%M)
NAME=io2websocket-gateway
PKG=${NAME}-${VERSION}-${BUILD}
DPKG_ARCH:=all
DPKG=dpkg-deb

all:

frontend:
	cd frontend && npm install
	cd frontend && npm run build
	cp frontend/dist/* public/

server: ${NAME}

dist: release/${PKG}.deb

build/${PKG}: server
	mkdir -p $@/etc/
	mkdir -p $@/usr/local/bin
	mkdir -p $@/usr/share/${NAME}/public
	mkdir -p $@/usr/lib/systemd/system/
	mkdir -p $@/etc/default/
	mkdir -p $@/var/lib/${NAME}
	mkdir -p $@/DEBIAN
	cp ${NAME} $@/usr/local/bin/
	cp config.txt $@/usr/share/${NAME}/io2websocket-gateway.conf
	cp -r frontend/dist/* $@/usr/share/${NAME}/public/
	cp dist/${NAME}.service $@/usr/lib/systemd/system/
	cp dist/postinst $@/DEBIAN/postinst
	chmod 755 $@/DEBIAN/postinst
	chmod a+x $@/usr/local/bin/${NAME}
	cp dist/dpkg-control $@/DEBIAN/control
	echo Version: ${VERSION}-${BUILD} >> $@/DEBIAN/control
	echo Architecture: ${DPKG_ARCH} >> $@/DEBIAN/control

release/${PKG}.deb: build/${PKG}
	mkdir -p release
	${DPKG} -b $< $@

clean:
	rm -rf build

${NAME}:
	go build
