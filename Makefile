# Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file.

# Build all by default, even if it's not first
.DEFAULT_GOAL := all

.PHONY: all
all: gen boilerplate lint test build

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/marmotedu/iam
VERSION_PACKAGE=github.com/marmotedu/component-base/pkg/version

# Copy githook scripts when execute makefile
COPY_GITHOOK:=$(shell cp -f githooks/* .git/hooks/)

# ==============================================================================
# Includes

include scripts/make-rules/common.mk
include scripts/make-rules/golang.mk
include scripts/make-rules/image.mk
include scripts/make-rules/deploy.mk
include scripts/make-rules/license.mk
include scripts/make-rules/gen.mk
include scripts/make-rules/ca.mk
include scripts/make-rules/release.mk
include scripts/make-rules/swagger.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG        Whether to generate debug symbols. Default is 0.
  BINS         The binaries to build. Default is all of cmd.
               This option is available when using: make build/build.multiarch
               Example: make build BINS="iam-apiserver iam-authz-server"
  IMAGES       Backend images to make. Default is all of cmd starting with iam-.
               This option is available when using: make image/image.multiarch/push/push.multiarch
               Example: make image.multiarch IMAGES="iam-apiserver iam-authz-server"
  PLATFORMS    The multiple platforms to build. Default is linux_amd64 and linux_arm64.
               This option is available when using: make build.multiarch/image.multiarch/push.multiarch
               Example: make image.multiarch IMAGES="iam-apiserver iam-pump" PLATFORMS="linux_amd64 linux_arm64"
  VERSION      The version information compiled into binaries.
               The default is obtained from git.
  V            Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets

## build: Build source code for host platform.
.PHONY: build
build:
	@$(MAKE) go.build

## build.multiarch: Build source code for multiple platforms. See option PLATFORMS.
.PHONY: build.multiarch
build.multiarch:
	@$(MAKE) go.build.multiarch

## image: Build docker images for host arch.
.PHONY: image
image:
	@$(MAKE) image.build

## image.multiarch: Build docker images for multiple platforms. See option PLATFORMS.
.PHONY: image.multiarch
image.multiarch:
	@$(MAKE) image.build.multiarch

## push: Build docker images for host arch and push images to registry.
.PHONY: push
push:
	@$(MAKE) image.push

## push.multiarch: Build docker images for multiple platforms and push images to registry.
.PHONY: push.multiarch
push.multiarch:
	@$(MAKE) image.push.multiarch

## deploy: Deploy updated components to development env.
.PHONY: deploy
deploy:
	@$(MAKE) deploy.run

## clean: Remove all files that are created by building.
.PHONY: clean
clean:
	@$(MAKE) go.clean

## lint: Check syntax and styling of go sources.
.PHONY: lint
lint: format
	@$(MAKE) go.lint

## test: Run unit test.
.PHONY: test
test:
	@$(MAKE) go.test

.PHONY: release.build
release.build:
	@$(MAKE) push.multiarch

## release: Release iam 
.PHONY: release
release:
	@$(MAKE) release.run

## format: Gofmt (reformat) package sources (exclude vendor dir if existed).
.PHONY: format 
format:
	@find . ! -path "./vendor/*" -name "*.go" | xargs gofmt -s -w

## boilerplate: Verify the boilerplate headers for all files.
.PHONY: boilerplate 
boilerplate:
	@$(MAKE) license.verify

## license: Ensures source code files have copyright license headers.
.PHONY: license 
license:
	@$(MAKE) license.add

## gen: Generate all necessary files, such as error code files.
.PHONY: gen
gen:
	@$(MAKE) gen.run

## ca: Generate CA files for all iam components.
.PHONY: ca
ca:
	@$(MAKE) gen.ca

## install: Install iam system with all its components.
.PHONY: install
install:
	@$(MAKE) install.install

## swagger: Generate swagger document.
.PHONY: swagger
swagger:
	@$(MAKE) swagger.run

## serve-swagger: Serve swagger spec and docs.
.PHONY: swagger.serve
serve-swagger:
	@$(MAKE) swagger.serve

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo -e "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"