include env

PROJECTNAME="flemzerd"
GO=vgo
PKGS=$(shell $(GO) list ./... | grep -v vendor)
VERSION=$(shell git describe --tags --always)
GOOS=$(shell $(GO) env GOOS)
GOARCH=$(shell $(GO) env GOARCH)

PACKAGE_NAME=flemzerd_$(VERSION)_$(GOOS)_$(GOARCH)

FLAGS=-X github.com/macarrie/flemzerd/configuration.Version=$(VERSION) -X github.com/macarrie/flemzerd/configuration.TRAKT_CLIENT_SECRET=$(FLZ_TRAKT_CLIENT_SECRET) -X github.com/macarrie/flemzerd/configuration.TELEGRAM_BOT_TOKEN=$(FLZ_TELEGRAM_BOT_TOKEN) -X github.com/macarrie/flemzerd/configuration.TMDB_API_KEY=$(FLZ_TMDB_API_KEY) -X github.com/macarrie/flemzerd/configuration.TVDB_API_KEY=$(FLZ_TVDB_API_KEY)

MAKEFLAGS += --silent

all: help
server/ui/node_modules:
	@echo " > Installing npm modules"
	@cd server/ui && npm install
	echo -e "\tNPM modules installed"

## webui: Build web interface and copy it into package folder
webui: server/ui/node_modules
	echo " > Building Web interface"
	echo -e "\tNode version: $$(node -v)"
	echo -e "\tNPM version: $$(npm -v)"
	mkdir -p ../../package/$(PACKAGE_NAME)/ui/
	cd server/ui && ./node_modules/@angular/cli/bin/ng build --prod --output-path "../../package/$(PACKAGE_NAME)/ui/"
	echo -e "\tInterface build complete: package/$(PACKAGE_NAME)/ui/"

tmp/vidocq/target/release/vidocq:
	echo " > Building dependencies: vidocq"
	-rm -rf tmp
	mkdir -p tmp
	git clone https://github.com/macarrie/vidocq tmp/vidocq
	cd tmp/vidocq && cargo build --release
	echo -e "\tVidocq build into tmp/vidocq/target/release/vidocq"

package/$(PACKAGE_NAME)/dependencies/vidocq: tmp/vidocq/target/release/vidocq
	mkdir -p package/$(PACKAGE_NAME)/dependencies
	cp tmp/vidocq/target/release/vidocq package/$(PACKAGE_NAME)/dependencies/vidocq
 
## bin: Build flemzerd binary
bin: 
	echo " > Building flemzerd binary"
	mkdir -p bin/
	$(GO) build -v -ldflags="$(FLAGS)" -o bin/flemzerd
	echo -e "\tBinary build complete: bin/flemzerd"

## package: Create package with binary, dependencies and installation scripts
package:
	echo " > Creating package"
	mkdir -p package/$(PACKAGE_NAME)
	cp -r install/* package/$(PACKAGE_NAME)
	cp -r bin package/$(PACKAGE_NAME)
	echo -e "\tPackage build complete: package/$(PACKAGE_NAME)"


## build: Build project (binary depencies, webui and package)
build: package/$(PACKAGE_NAME)/dependencies/vidocq webui bin package

## doc: Build documentation
doc:
	echo " > Building documentation files"
	cp package/$(PACKAGE_NAME)/ui/styles*.css docs_src/themes/flemzer/static/css/flemzer.css

## install: Install flemzerd on your machine
install: build
	echo " > Launching installation script"
	cd package/$(PACKAGE_NAME)/ && sudo ./install.sh && cd ..

## update: Update flemzerd on your machine
update: build
	echo " > Launching update script"
	cd package/$(PACKAGE_NAME)/ && sudo ./update.sh && cd ..

## test: Run unit tests
test:
	echo " > Launching unit tests"
	-rm -rf cover
	tests_failed=0
	mkdir -p cover
	echo "mode: count" > cover/coverage.cov
	for d in $(PKGS); \
	do \
		tests_in_package=$$(ls $$GOPATH/src/$$d | grep _test.go | wc -l); \
		if [ $$tests_in_package -gt 0 ]; \
		then \
			$(GO) test -covermode=count -coverprofile "cover/$${d##*/}.cov" "$$d"; \
			ret=$$?; \
			if [ $$ret -ne 0 ]; \
			then \
				tests_failed=$$ret; \
			fi; \
		fi; \
	done;
	tail -q -n +2 cover/*.cov >> cover/coverage.cov
	$(GO) tool cover -func=cover/coverage.cov
	exit $$tests_failed;


## watch: Run command when code changes (make watch run="go test ./...")
watch:
	inotifywait -m -r ./ -e modify,create,delete --exclude '[^g][^o]$$' --format '%T' --timefmt '%s' | \
	stdbuf -oL uniq | \
	while read f; \
		echo -e "\n\n==================== Change detected ===================="; \
		do $(run); \
	done

start-server: stop-server
	touch /tmp/flemzerd.pid
	sudo -E -u flemzer ./bin/flemzerd -d & echo $$! > /tmp/flemzerd.pid

stop-server:
	-kill $$(cat /tmp/flemzerd.pid)
	-rm /tmp/flemzerd.pid

## start: Start dev server. Change in go files restarts compilation and server
start:
	echo " > Starting dev server"
	@bash -c "trap 'make stop' EXIT; $(MAKE) bin start-server watch run='make bin start-server'"

## stop: Stop dev server
stop: stop-server

## reload: Reload configuration by sending USR1 signal to flemzerd daemon
reload:
	-kill -USR1 $$(cat /tmp/flemzerd.pid)


## clean: Clean build and tests artifacts
clean:
	echo " > Cleaning build artifacts"
	-rm -rf package
	-rm -rf cover
	-rm -rf tmp
	-rm -rf bin

.PHONY: help webui bin build doc test install update clean package
help: Makefile
	echo " Choose a command run in "$(PROJECTNAME)":"
	sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'

