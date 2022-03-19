LDFLAGS=-w -s
PKGVERSION=$(shell cat main.go | grep VERSION | awk '{print $$5}' | sed 's/"//g')
PKGDATE=$(shell date -R)

.PHONY: \
	build \
	make-deb \
	plugins \
	clean \
	all
	linux-arm \
	linux-arm64 \
	linux-amd64 \
	windows-amd64 \
	darwin-amd64 \
	darwin-arm64

linux-amd64:
	GOOS=linux GOARCH=amd64 EXT= $(MAKE) build
	GOOS=linux GOARCH=amd64 EXT=.so $(MAKE) plugins
	GOOS=linux GOARCH=amd64 ${MAKE} make-deb

linux-arm:
	GOOS=linux GOARCH=arm EXT= $(MAKE) build

linux-arm64:
	GOOS=linux GOARCH=arm64 EXT= $(MAKE) build
	GOOS=linux GOARCH=arm64 EXT= $(MAKE) make-deb

windows-amd64:
	GOOS=windows GOARCH=amd64 EXT=.exe $(MAKE) build

darwin-amd64:
	GOOS=darwin GOARCH=amd64 EXT= ${MAKE} build

darwin-arm64:
	GOOS=darwin GOARCH=arm64 EXT= ${MAKE} build

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="$(LDFLAGS)" -o BUILD/$(GOOS)-$(GOARCH)/exequte$(EXT) main.go
	cp exequte.json BUILD/$(GOOS)-$(GOARCH)

plugins:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="${LDFLAGS}" -buildmode=plugin -o BUILD/$(GOOS)-$(GOARCH)/plugins/dummy$(EXT) plugins/dummy/dummy.go

make-deb:
	mkdir -p BUILD/debian-$(GOARCH)/debian/exequte
	
	cp BUILD/$(GOOS)-$(GOARCH)/exequte BUILD/debian-$(GOARCH)/exequte
	cp exequte.json BUILD/debian-$(GOARCH)/exequte.json
	cp install/systemd/exequte.service BUILD/debian-$(GOARCH)/debian/exequte.service

	cp install/deb/compat BUILD/debian-$(GOARCH)/debian/compat
	cp install/deb/copyright BUILD/debian-$(GOARCH)/debian/copyright
	cp install/deb/postinst BUILD/debian-$(GOARCH)/debian/postinst
	cp install/deb/rules BUILD/debian-$(GOARCH)/debian/rules
	
	cat install/deb/changelog | sed 's/__PKGDATE__/$(PKGDATE)/' | sed s'/__PKGVERSION__/$(PKGVERSION)/' > BUILD/debian-$(GOARCH)/debian/changelog
	cat install/deb/control | sed 's/__PKGARCH__/$(GOARCH)/' > BUILD/debian-$(GOARCH)/debian/control

	cd BUILD/debian-$(GOARCH); dpkg-buildpackage -us -uc -b -a $(GOARCH)
	lintian --check --color always BUILD/*_$(GOARCH).deb
	
clean:
	rm -rf BUILD/debian-amd64
	rm -rf BUILD/debian-arm64
	rm -rf BUILD/linux-arm
	rm -rf BUILD/linux-arm64
	rm -rf BUILD/linux-amd64
	rm -rf BUILD/windows-amd64
	rm -rf BUILD/darwin-amd64
	rm -rf BUILD/darwin-arm64
	rm -rf BUILD/*.deb
	rm -rf BUILD/*.changes
	rm -rf BUILD/*.buildinfo

all: linux-arm linux-arm64 linux-amd64 windows-amd64 darwin-amd64 darwin-arm64
