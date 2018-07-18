PKGS=$(shell vgo list ./... | grep -v vendor)

VERSION=$(shell git describe --tags --always)
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
GOPATH=~/go

PACKAGE_NAME=flemzerd_$(VERSION)_$(GOOS)_$(GOARCH)

FLAGS=-X github.com/macarrie/flemzerd/configuration.Version=$(VERSION) -X github.com/macarrie/flemzerd/configuration.TRAKT_CLIENT_SECRET=$(FLZ_TRAKT_CLIENT_SECRET) -X github.com/macarrie/flemzerd/configuration.TELEGRAM_BOT_TOKEN=$(FLZ_TELEGRAM_BOT_TOKEN) -X github.com/macarrie/flemzerd/configuration.TMDB_API_KEY=$(FLZ_TMDB_API_KEY) -X github.com/macarrie/flemzerd/configuration.TVDB_API_KEY=$(FLZ_TVDB_API_KEY)


all: build

server/ui/node_modules:
	cd server/ui && npm install

webui: server/ui/node_modules
	node -v
	npm -v
	mkdir -p ../../package/$(PACKAGE_NAME)/ui/
	cd server/ui && ./node_modules/@angular/cli/bin/ng build --prod --output-path "../../package/$(PACKAGE_NAME)/ui/"

tmp/vidocq/target/release/vidocq:
	-rm -rf tmp
	mkdir -p tmp
	git clone https://github.com/macarrie/vidocq tmp/vidocq
	cd tmp/vidocq && cargo build --release

package/$(PACKAGE_NAME)/dependencies/vidocq: tmp/vidocq/target/release/vidocq
	mkdir -p package/$(PACKAGE_NAME)/dependencies
	cp tmp/vidocq/target/release/vidocq package/$(PACKAGE_NAME)/dependencies/vidocq
 
bin: 
	mkdir -p bin/
	vgo build -v -ldflags="$(FLAGS)" -o bin/flemzerd

package:
	mkdir -p package/$(PACKAGE_NAME)
	cp -r install/* package/flemzerd_$(VERSION)_$(GOOS)_$(GOARCH)
	cp -r bin package/flemzerd_$(VERSION)_$(GOOS)_$(GOARCH)


build: package/$(PACKAGE_NAME)/dependencies/vidocq webui bin package

#doc: webui
	#cp server/ui/css/flemzer.css docs_src/themes/flemzer/static/css/flemzer.css

install: build
	cd package/$(PACKAGE_NAME)/install/ && sudo ./install.sh && cd ..

update: build
	cd package/$(PACKAGE_NAME)/install/ && sudo ./update.sh && cd ..

test:
	-rm -rf cover
	@tests_failed=0
	@mkdir -p cover
	@echo "mode: count" > cover/coverage.cov
	@for d in $(PKGS); \
	do \
		tests_in_package=$$(ls $(GOPATH)/src/$$d | grep _test.go | wc -l); \
		if [ $$tests_in_package -gt 0 ]; \
		then \
			vgo test -covermode=count -coverprofile "cover/$${d##*/}.cov" "$$d"; \
			ret=$$?; \
			if [ $$ret -ne 0 ]; \
			then \
				tests_failed=$$ret; \
			fi; \
		fi; \
	done;
	@tail -q -n +2 cover/*.cov >> cover/coverage.cov
	@vgo tool cover -func=cover/coverage.cov
	@exit $$tests_failed;


clean:
	-rm flemzerd
	-rm -rf package
	-rm -rf cover
	-rm -rf tmp

.PHONY: all webui bin build doc test install update clean package
