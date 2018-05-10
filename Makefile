BINARY=WhyAwake.app/Contents/MacOS/whyawake
SOURCEDIR=.
LIBDIR=../menuet ../go-assertions ../go-caffeinate
SOURCES := $(shell find $(SOURCEDIR) $(LIBDIR) -name '*.go' -o -name '*.m' -o -name '*.h' -o -name '*.c') Makefile

run: $(BINARY)
	./$(BINARY)

$(BINARY): $(SOURCES)
	@echo $(SOURCES)
	go build -o $(BINARY)

clean:
	rm -f $(BINARY)