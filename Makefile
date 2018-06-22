PKGS=$(shell vgo list ./... | grep -v vendor)
VERSION=$(shell git describe --tags --always)
FLAGS=-X main.version=$(VERSION)

all: build

webui:
	@sass --version
	#sass server/ui/css/flemzer.scss server/ui/css/flemzer.css

install/vidocq:
	-rm -rf tmp
	mkdir -p tmp
	git clone https://github.com/macarrie/vidocq tmp/vidocq
	cd tmp/vidocq && cargo build --release
	cp tmp/vidocq/target/release/vidocq install/vidocq

build: install/vidocq webui
	vgo build -v -ldflags="$(FLAGS)"

doc: webui
	cp server/ui/css/flemzer.css docs_src/themes/flemzer/static/css/flemzer.css

install: build
	cd install && sudo ./install.sh && cd ..

update: build
	cd install && sudo ./update.sh && cd ..

test:
	-rm -rf cover
	@tests_failed=0
	@mkdir -p cover
	@echo "mode: count" > cover/coverage.cov
	@for d in $(PKGS); \
	do \
		vgo test -covermode=count -coverprofile "cover/$${d##*/}.cov" "$$d"; \
		ret=$$?; \
		if [ $$ret -ne 0 ]; \
		then \
			tests_failed=$$ret; \
		fi; \
	done;
	@tail -q -n +2 cover/*.cov >> cover/coverage.cov
	#@vgo tool cover -func=cover/coverage.cov
	@exit $$tests_failed;


clean:
	-rm flemzerd
	-rm coverage.txt
	-rm server/ui/css/*.css server/ui/css/*.map
	-rm -rf .sass-cache
	-rm -rf tmp
	-rm -rf install/vidocq

.PHONY: all webui deps build doc test install update clean
