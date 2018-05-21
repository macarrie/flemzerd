all: build

webui:
	sass server/ui/css/flemzer.scss server/ui/css/flemzer.css

build: webui
	go build -v

doc: webui
	cp server/ui/css/flemzer.css docs_src/themes/flemzer/static/css/flemzer.css

test:
	./test.sh

clean:
	-rm flemzerd
	-rm server/ui/css/*.css server/ui/css/*.map
	-rm -rf .sass-cache
