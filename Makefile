.PHONY: build clean run

GO ?= go

SRC_DIR = .
DIST_DIR = ./build

BINARY = vserial
ifeq (${GOOS}, windows)
    BINARY := $(BINARY).exe
endif

BUILD_ARGS = -v -trimpath
BUILD_FLAGS = -s -w

build:
	@echo "[Info] Building project, output file path: $(DIST_DIR)/$(BINARY)"
	@mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GOARM=${GOARM} GOMIPS=${GOMIPS} \
		$(GO) build -ldflags="$(BUILD_FLAGS)" $(BUILD_ARGS) -o $(DIST_DIR)/$(BINARY) $(SRC_DIR)
	@echo "[Info] Build completed."

run:
	@mkdir -p $(DIST_DIR)
	@echo "[Info] Running project..."
	$(GO) run -gcflags="all=-N -l" -race $(SRC_DIR)/*.go -database $(DIST_DIR)/states.db.local

clean:
	@echo "[Warn] Cleaning up project..."
	@rm -rf $(DIST_DIR)/*
