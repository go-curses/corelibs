#!/usr/bin/env make

SHELL=/bin/bash

LOCAL_CDK_PATH ?= ../../cdk
LOCAL_CTK_PATH ?= ../../ctk

.PHONY: help local unlocal be-update tidy

help:
	@echo "# make <local|unlocal|be-update|tidy>"

local:
	@for dir in `find . -name "go.mod" -exec dirname \{\} \;`; do \
		pushd "$${dir}" > /dev/null; \
		if egrep -q 'go-curses/cdk v' go.mod; then \
			echo "# go mod local cdk: $${dir}"; \
			go mod edit -replace=github.com/go-curses/cdk=${LOCAL_CDK_PATH}; \
		fi; \
		if egrep -q 'go-curses/ctk v' go.mod; then \
			echo "# go mod local ctk: $${dir}"; \
			go mod edit -replace=github.com/go-curses/ctk=${LOCAL_CTK_PATH}; \
		fi; \
		popd > /dev/null; \
	done

unlocal:
	@for dir in `find . -name "go.mod" -exec dirname \{\} \;`; do \
		pushd "$${dir}" > /dev/null; \
		if egrep -q 'go-curses/cdk v' go.mod; then \
			echo "# go mod unlocal cdk: $${dir}"; \
			go mod edit -dropreplace=github.com/go-curses/cdk; \
		fi; \
		if egrep -q 'go-curses/ctk v' go.mod; then \
			echo "# go mod unlocal ctk: $${dir}"; \
			go mod edit -dropreplace=github.com/go-curses/ctk; \
		fi; \
		popd > /dev/null; \
	done

be-update:
	@for dir in `find . -name "go.mod" -exec dirname \{\} \;`; do \
		pushd "$${dir}" > /dev/null; \
		if egrep -q 'go-curses/cdk v' go.mod; then \
			echo "# go get cdk: $${dir}"; \
			go get github.com/go-curses/cdk@latest; \
		fi; \
		if egrep -q 'go-curses/ctk v' go.mod; then \
			echo "# go get ctk: $${dir}"; \
			go get github.com/go-curses/ctk@latest; \
		fi; \
		popd > /dev/null; \
	done

tidy:
	@for dir in `find . -name "go.mod" -exec dirname \{\} \;`; do \
		pushd "$${dir}" > /dev/null; \
		echo "# go mod tidy: $${dir}"; \
		go mod tidy; \
		popd > /dev/null; \
	done
