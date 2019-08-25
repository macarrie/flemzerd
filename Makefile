include env

PROJECTNAME="flemzerd"
GO=go
PKGS=$(shell $(GO) list ./... | grep -v vendor | sed 's\#github.com/macarrie/flemzerd/*\#./\#')
VERSION=$(shell git describe --tags --always)
GOOS=$(shell $(GO) env GOOS)
GOARCH=$(shell $(GO) env GOARCH)

CARGO=$(shell which cargo)

PACKAGE_NAME=flemzerd_$(VERSION)_$(GOOS)_$(GOARCH)
VIDOCQ_VERSION=v0.1.1

LDFLAGS=-X github.com/macarrie/flemzerd/configuration.Version=$(VERSION) -X github.com/macarrie/flemzerd/configuration.TRAKT_CLIENT_SECRET=$(FLZ_TRAKT_CLIENT_SECRET) -X github.com/macarrie/flemzerd/configuration.TELEGRAM_BOT_TOKEN=$(FLZ_TELEGRAM_BOT_TOKEN) -X github.com/macarrie/flemzerd/configuration.TMDB_API_KEY=$(FLZ_TMDB_API_KEY) -X github.com/macarrie/flemzerd/configuration.TVDB_API_KEY=$(FLZ_TVDB_API_KEY)

MAKEFLAGS += --silent

all: help
server/ui/node_modules:
	@echo " > Installing npm modules"
	@cd server/ui && npm --unsafe-perm install
	echo -e "\tNPM modules installed"

## webui: Build web interface and copy it into package folder
webui: server/ui/node_modules
	echo " > Building Web interface"
	echo -e "\tNode version: $$(node -v)"
	echo -e "\tNPM version: $$(npm -v)"
	mkdir -p package/$(PACKAGE_NAME)/ui/
	cd server/ui && npm run build
	cp -r server/ui/build/* package/$(PACKAGE_NAME)/ui/
	echo -e "\tInterface build complete: package/$(PACKAGE_NAME)/ui/"

tmp/vidocq:
	echo " > Building dependencies: vidocq"
	-rm -rf tmp
	mkdir -p tmp
	curl -L https://github.com/macarrie/vidocq/releases/download/$(VIDOCQ_VERSION)/vidocq -o tmp/vidocq
	echo -e "\tVidocq executable downloaded into tmp/vidocq"

package/$(PACKAGE_NAME)/dependencies/vidocq: tmp/vidocq
	mkdir -p package/$(PACKAGE_NAME)/dependencies
	cp tmp/vidocq package/$(PACKAGE_NAME)/dependencies/vidocq

## bin: Build flemzerd binary
bin:
	echo " > Building flemzerd binary"
	mkdir -p bin/
	CC=gcc $(GO) build -v -ldflags="$(LDFLAGS)" -o bin/flemzerd
	echo -e "\tBinary build complete: bin/flemzerd"


## docker: Build flemzerd docker image
docker:
	echo " > Building docker image"
	docker build -t macarrie/flemzerd .
	echo -e "\tDocker build complete"


## package: Create package with binary, dependencies and installation scripts
package:
	echo " > Creating package"
	mkdir -p package/$(PACKAGE_NAME)
	cp -r install/* package/$(PACKAGE_NAME)
	cp -r bin package/$(PACKAGE_NAME)
	sed 's#build: .#image: macarrie/flemzerd#' docker-compose.yml > package/$(PACKAGE_NAME)/docker-compose.yml
	echo -e "\tPackage build complete: package/$(PACKAGE_NAME)"


## build: Build project (binary depencies, webui and package)
build: package/$(PACKAGE_NAME)/dependencies/vidocq webui bin package


node_modules/node-sass/bin/node-sass:
	npm install node-sass


## doc: Build documentation
doc: node_modules/node-sass/bin/node-sass
	echo " > Building documentation files"
	./node_modules/node-sass/bin/node-sass --output-style compressed server/ui/src/styles.scss docs_src/themes/flemzer/static/css/flemzer.css

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
		tests_in_package=$$(ls ./$$d | grep _test.go | wc -l); \
		if [ $$tests_in_package -gt 0 ]; \
		then \
			$(GO) test -covermode=count -coverprofile "cover/$${d##*/}.testcov" "$$d"; \
			ret=$$?; \
			if [ $$ret -ne 0 ]; \
			then \
				exit $$ret; \
			fi; \
		fi; \
	done;
	tail -q -n +2 cover/*.testcov >> cover/coverage.cov
	$(GO) tool cover -func=cover/coverage.cov


## watch: Run command when code changes (make watch run="go test ./...")
watch:
	inotifywait -m -r ./ -e modify,create,delete --exclude '[^g][^o]$$' --format '%T' --timefmt '%s' | \
	stdbuf -oL uniq | \
	while read f; \
		echo -e "\n\n==================== Change detected ===================="; \
		do $(run); \
	done

#start-server: stop-server
start-server:
	#sudo -E -u flemzer ./bin/flemzerd -d & echo $$! > /tmp/flemzerd.pid
	sudo -E ./bin/flemzerd -d

stop-server:
	-kill $$(cat /tmp/flemzerd.pid)
	-rm /tmp/flemzerd.pid

## start: Start dev server. Change in go files restarts compilation and server
start:
	echo " > Starting dev server"
	$(MAKE) bin start-server

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

