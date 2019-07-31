PREFIX=/usr/local
PKG := gitlab.com/gitlab-org/gitlab-workhorse
BUILD_DIR ?= $(CURDIR)
TARGET_DIR ?= $(BUILD_DIR)/_build
TARGET_SETUP := $(TARGET_DIR)/.ok
BIN_BUILD_DIR := $(TARGET_DIR)/bin
TOOLS_DIR := $(BUILD_DIR)/_tools
PKG_BUILD_DIR := $(TARGET_DIR)/src/$(PKG)
COVERAGE_DIR := $(TARGET_DIR)/cover
VERSION_STRING := $(shell git describe)
ifeq ($(strip $(VERSION_STRING)),)
VERSION_STRING := v$(shell cat VERSION)
endif
VERSION := ${VERSION_STRING}-$(shell date -u +%Y%m%d.%H%M%S)
GOBUILD := go build -ldflags "-X main.Version=$(VERSION)"
EXE_ALL := gitlab-zip-cat gitlab-zip-metadata gitlab-workhorse
INSTALL := install
BUILD_TAGS := tracer_static tracer_static_jaeger
STATICCHECK_VERSION := 2019.2.2
STATICCHECK := $(TOOLS_DIR)/staticcheck
STATICCHECK_VERSION_INSTALLED = $(shell $(STATICCHECK) -version)
MINIMUM_SUPPORTED_GO_VERSION := 1.8

# Some users may have these variables set in their environment, but doing so could break
# their build process, so unset then
unexport GOROOT
unexport GOBIN

export GOPATH := $(TARGET_DIR)
export PATH := $(GOPATH)/bin:$(PATH)

# Returns a list of all non-vendored (local packages)
LOCAL_PACKAGES = $(shell cd "$(PKG_BUILD_DIR)" && GOPATH=$(GOPATH) go list ./... | grep -v -e '^$(PKG)/vendor/' -e '^$(PKG)/ruby/')
LOCAL_GO_FILES = $(shell find -L $(PKG_BUILD_DIR)  -name "*.go" -not -path "$(PKG_BUILD_DIR)/vendor/*" -not -path "$(PKG_BUILD_DIR)/_build/*")

define message
	@echo "### $(1)"
endef

.NOTPARALLEL:

.PHONY:	all
all:	clean-build $(EXE_ALL)

$(TARGET_SETUP):
	$(call message,"Setting up target directory")
	rm -rf $(TARGET_DIR)
	mkdir -p "$(dir $(PKG_BUILD_DIR))"
	ln -sf "$(CURDIR)" "$(PKG_BUILD_DIR)"
	mkdir -p "$(BIN_BUILD_DIR)"
	touch "$(TARGET_SETUP)"

gitlab-zip-cat:	$(TARGET_SETUP) $(shell find cmd/gitlab-zip-cat/ -name '*.go')
	$(call message,Building $@)
	$(GOBUILD) -tags "$(BUILD_TAGS)" -o $(BUILD_DIR)/$@ $(PKG)/cmd/$@

gitlab-zip-metadata:	$(TARGET_SETUP) $(shell find cmd/gitlab-zip-metadata/ -name '*.go')
	$(call message,Building $@)
	$(GOBUILD) -tags "$(BUILD_TAGS)" -o $(BUILD_DIR)/$@ $(PKG)/cmd/$@

gitlab-workhorse:	$(TARGET_SETUP) $(shell find . -name '*.go' | grep -v '^\./_')
	$(call message,Building $@)
	$(GOBUILD) -tags "$(BUILD_TAGS)" -o $(BUILD_DIR)/$@ $(PKG)

.PHONY:	install
install:	gitlab-workhorse gitlab-zip-cat gitlab-zip-metadata
	$(call message,$@)
	mkdir -p $(DESTDIR)$(PREFIX)/bin/
	cd $(BUILD_DIR) && $(INSTALL) gitlab-workhorse gitlab-zip-cat gitlab-zip-metadata $(DESTDIR)$(PREFIX)/bin/

.PHONY:	test
test: $(TARGET_SETUP) prepare-tests
	$(call message,$@)
	@go test -tags "$(BUILD_TAGS)" $(LOCAL_PACKAGES)
	@echo SUCCESS

.PHONY:	coverage
coverage:	$(TARGET_SETUP) prepare-tests
	$(call message,$@)
	@go test -tags "$(BUILD_TAGS)" -cover -coverprofile=test.coverage $(LOCAL_PACKAGES)
	go tool cover -html=test.coverage -o coverage.html
	rm -f test.coverage

.PHONY:	clean
clean:	clean-workhorse clean-build
	$(call message,$@)
	rm -rf testdata/data testdata/scratch

.PHONY:	clean-workhorse
clean-workhorse:
	$(call message,$@)
	rm -f $(EXE_ALL)

.PHONY:	tag
tag:
	$(call message,$@)
	sh _support/tag.sh

.PHONY:	signed_tag
signed_tag:
	$(call message,$@)
	TAG_OPTS=-s sh _support/tag.sh

.PHONY:	clean-build
clean-build:
	$(call message,$@)
	rm -rf $(TARGET_DIR)

.PHONY:	prepare-tests
prepare-tests:	govendor-sync testdata/data/group/test.git $(EXE_ALL)

testdata/data/group/test.git:
	$(call message,$@)
	git clone --quiet --bare https://gitlab.com/gitlab-org/gitlab-test.git $@

.PHONY: verify
verify: lint vet detect-context check-formatting staticcheck

.PHONY: lint
lint: $(TARGET_SETUP) govendor-sync
	$(call message,Verify: $@)
	@command -v golint || go get -v golang.org/x/lint/golint
	@_support/lint.sh $(LOCAL_PACKAGES)

.PHONY: vet
vet: $(TARGET_SETUP) govendor-sync
	$(call message,Verify: $@)
	@go vet $(LOCAL_PACKAGES)

.PHONY: detect-context
detect-context: $(TARGET_SETUP)
	$(call message,Verify: $@)
	_support/detect-context.sh

.PHONY: check-formatting
check-formatting: $(TARGET_SETUP) install-goimports
	$(call message,Verify: $@)
	@_support/validate-formatting.sh $(LOCAL_GO_FILES)

$(STATICCHECK):
	mkdir -p $(TOOLS_DIR)
	_support/install-staticcheck.sh $(STATICCHECK_VERSION) $@

# staticcheck will tailor some responses given a minimum Go version, so pass that through the CLI
.PHONY: staticcheck
staticcheck: $(TARGET_SETUP) govendor-sync $(STATICCHECK)
	$(call message,Verify: $@)
	@if [ "staticcheck $(STATICCHECK_VERSION)" != "$(STATICCHECK_VERSION_INSTALLED)" ]; then \
		echo wrong version of staticcheck. required "staticcheck $(STATICCHECK_VERSION)", you have "$(STATICCHECK_VERSION_INSTALLED)". \
		echo make clean-tools and try again.\
		exit 1; \
	fi
	@$(TOOLS_DIR)/staticcheck -go $(MINIMUM_SUPPORTED_GO_VERSION) $(LOCAL_PACKAGES)

# Some vendor components, used for testing are GPL, so we don't distribute them
# and need to go a sync before using them
.PHONY: govendor-sync
govendor-sync: $(TARGET_SETUP)
	$(call message,$@)
	@command -v govendor || go get github.com/kardianos/govendor
	@cd $(PKG_BUILD_DIR) && govendor sync

# In addition to fixing imports, goimports also formats your code in the same style as gofmt
# so it can be used as a replacement.
.PHONY: fmt
fmt: $(TARGET_SETUP) install-goimports
	$(call message,$@)
	@goimports -w -local $(PKG) -l $(LOCAL_GO_FILES)

.PHONY:	goimports
install-goimports:	$(TARGET_SETUP)
	$(call message,$@)
	@command -v goimports || go get -v golang.org/x/tools/cmd/goimports

.PHONY:	clean-tools
clean-tools:
	$(call message,$@)
	rm -rf $(TOOLS_DIR)
