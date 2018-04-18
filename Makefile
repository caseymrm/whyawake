BINARY=cantsleep
SOURCEDIR=.
LIBDIR=../go-statusbar/tray ../go-assertions
SOURCES := $(shell find $(SOURCEDIR) $(LIBDIR) -name '*.go' -o -name '*.m' -o -name '*.h' -o -name '*.c') Makefile

run: $(BINARY)
	./$(BINARY)

$(BINARY): $(SOURCES)
	@echo $(SOURCES)
	go build -o $(BINARY)
	cp $(BINARY) CantSleep.app/Contents/MacOS/

clean:
	rm -f $(BINARY)