all: build

webui:
	sass server/ui/css/flemzer.scss server/ui/css/flemzer.css

build: webui
	go build -v

test:
	./test.sh

clean:
	-rm flemzerd
	-rm server/ui/css/*.css server/ui/css/*.map
	-rm -rf .sass-cache
