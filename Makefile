build: clean
	go build -o bin/gorefs gorefs.go

clean:
	-rm -rf bin

install: uninstall build 
	cp bin/gorefs ~/.local/bin/gorefs

uninstall:
	-rm -f ~/.local/bin/gorefs