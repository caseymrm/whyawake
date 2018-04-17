BINARY=cantsleep
SOURCEDIR=.
LIBDIR=../go-statusbar/tray/
SOURCES := $(shell find $(SOURCEDIR) $(LIBDIR) -name '*.go' -o -name '*.m' -o -name '*.h' -o -name '*.c')

run: $(BINARY)
	./$(BINARY)

$(BINARY): $(SOURCES)
	go build -o $(BINARY)
	cp $(BINARY) CantSleep.app/Contents/MacOS/

clean:
	rm -f $(BINARY)