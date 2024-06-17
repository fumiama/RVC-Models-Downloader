PROJECT_NAME := rvcmd
BUILD_PATH := build
RUN_CMD := -h

GOOS := $(go env GOOS)
GOARCH := $(go env GOARCH)
GO_SRC_WIN := $(shell find . -maxdepth 1 -name '*.go' | grep -v '_test.go$$' | grep -v '_windows.go$$')
GO_SRC_NO_WIN := $(shell echo $(GO_SRC_WIN) | grep -v '_windows.go$$')

all:
	@$(MAKE) -e bin
bin: gen dir tidy
	@if [[ "$(GOARCH)" == "windows" ]]; then \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w" -trimpath -o $(BUILD_PATH)/$(PROJECT_NAME).exe; \
	else \
		GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags "-s -w" -trimpath -o $(BUILD_PATH)/$(PROJECT_NAME); \
	fi
run: bin
	@if [[ "$(GOARCH)" == "windows" ]]; then \
		$(BUILD_PATH)/$(PROJECT_NAME).exe $(RUN_CMD); \
	else \
		$(BUILD_PATH)/$(PROJECT_NAME) $(RUN_CMD); \
	fi
gen:
	@go generate
tidy:
	@go mod tidy
dir:
	@if [ ! -d "$(BUILD_PATH)" ]; then mkdir $(BUILD_PATH); fi
clean:
	@if [ -d "$(BUILD_PATH)" ]; then \
		rm -rf $(BUILD_PATH)/$(PROJECT_NAME)*; \
	fi
