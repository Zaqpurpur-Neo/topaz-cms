GO_BIN   := topaz-server
VALA_BIN := topaz-gui

GO        := go
VALAC     := valac
WEBKIT_PKG := $(shell pkg-config --exists webkit2gtk-4.1 && echo "webkit2gtk-4.1" || echo "webkit2gtk-4.0")

VALAFLAGS := --pkg gtk+-3.0 --pkg gio-2.0 --pkg $(WEBKIT_PKG)

.PHONY: all build clean run check-deps

all: build

check-deps:
	@which $(VALAC) > /dev/null 2>&1 || (echo "[ERROR]: valac is missing. Install with: sudo pacman -S vala" && exit 1)
	@pkg-config --exists gtk+-3.0 || (echo "[ERROR]: gtk+-3.0 missing. Install with: sudo pacman -S gtk3" && exit 1)
	@pkg-config --exists gio-2.0 || (echo "[ERROR]: gio-2.0 missing." && exit 1)
	@pkg-config --exists $(WEBKIT_PKG) || (echo "[ERROR]: webkit2gtk missing. Install with: sudo pacman -S webkit2gtk-4.1" && exit 1)

build: check-deps $(GO_BIN) $(VALA_BIN)

$(GO_BIN): main.go
	@echo "[BUILD]: Compile GO"
	CGO_ENABLED=0 $(GO) build -ldflags="-s -w" -o $(GO_BIN) main.go

$(VALA_BIN): $(wildcard *.vala)
	@echo "[BUILD]: Compile Vala"
	$(VALAC) $(VALAFLAGS) $^ -o $(VALA_BIN)

run: build
	@echo "[RUN]: HHH"
	./$(VALA_BIN)

clean:
	@echo "[CLEAN]: BBB"
	rm -f $(GO_BIN) $(VALA_BIN)
