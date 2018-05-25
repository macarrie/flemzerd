PKGS := $(shell go list ./... | grep -v vendor)

all: build

webui:
	sass server/ui/css/flemzer.scss server/ui/css/flemzer.css

pull:
	git pull

build: webui
	vgo build -v

doc: webui
	cp server/ui/css/flemzer.css docs_src/themes/flemzer/static/css/flemzer.css

install: pull build
	cd install && sudo ./install.sh && cd ..

update: pull build
	cd install && sudo ./update.sh && cd ..

test:
	@echo "" > coverage.txt
	@for d in $(PKGS); \
	do \
		go test -race -coverprofile=profile.out -covermode=atomic "$$d" ;\
		if [ -f profile.out ]; \
		then \
			cat profile.out >> coverage.txt ; \
			rm profile.out ; \
		fi \
	done \

clean:
	-rm flemzerd
	-rm coverage.txt
	-rm server/ui/css/*.css server/ui/css/*.map
	-rm -rf .sass-cache

.PHONY: all webui pull deps build doc test install update clean
