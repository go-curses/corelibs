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
			echo "# $${dir}: go mod local cdk"; \
			go mod edit -replace=github.com/go-curses/cdk=${LOCAL_CDK_PATH}; \
		fi; \
		if egrep -q 'go-curses/ctk v' go.mod; then \
			echo "# $${dir}: go mod local ctk"; \
			go mod edit -replace=github.com/go-curses/ctk=${LOCAL_CTK_PATH}; \
		fi; \
		FOUND_LIBS=`grep -h -v '^module' go.mod | grep 'go-curses/corelibs/' | perl -pe 's!^.*(github\.com/go-curses/corelibs/\S+).*$$!$${1}!'`; \
		if [ -n "$${FOUND_LIBS}" ]; then \
			for found_lib in $${FOUND_LIBS}; do \
				name=`basename $${found_lib}`; \
				echo "# $${dir}: go mod local ../$${name}"; \
				go mod edit -replace=$${found_lib}=../$${name}; \
			done; \
		fi; \
		popd > /dev/null; \
	done

unlocal:
	@for dir in `find . -name "go.mod" -exec dirname \{\} \;`; do \
		pushd "$${dir}" > /dev/null; \
		if egrep -q 'go-curses/cdk v' go.mod; then \
			echo "# $${dir}: go mod unlocal cdk"; \
			go mod edit -dropreplace=github.com/go-curses/cdk; \
		fi; \
		if egrep -q 'go-curses/ctk v' go.mod; then \
			echo "# $${dir}: go mod unlocal ctk"; \
			go mod edit -dropreplace=github.com/go-curses/ctk; \
		fi; \
		FOUND_LIBS=`grep -h -v '^module' go.mod | grep 'go-curses/corelibs/' | perl -pe 's!^.*(github\.com/go-curses/corelibs/\S+).*$$!$${1}!'`; \
		if [ -n "$${FOUND_LIBS}" ]; then \
			for found_lib in $${FOUND_LIBS}; do \
				name=`basename $${found_lib}`; \
				echo "# $${dir}: go mod unlocal corelibs/$${name}"; \
				go mod edit -dropreplace=$${found_lib}; \
			done; \
		fi; \
		popd > /dev/null; \
	done

be-update: export GOPROXY=direct
be-update:
	@for dir in `find . -name "go.mod" -exec dirname \{\} \;`; do \
		pushd "$${dir}" > /dev/null; \
		if egrep -q 'go-curses/cdk v' go.mod; then \
			echo "# $${dir}: go get cdk"; \
			go get github.com/go-curses/cdk@latest; \
		fi; \
		if egrep -q 'go-curses/ctk v' go.mod; then \
			echo "# $${dir}: go get ctk"; \
			go get github.com/go-curses/ctk@latest; \
		fi; \
		FOUND_LIBS=`grep -h -v '^module' go.mod | grep 'go-curses/corelibs/' | perl -pe 's!^.*(github\.com/go-curses/corelibs/\S+).*$$!$${1}!'`; \
		if [ -n "$${FOUND_LIBS}" ]; then \
			for found_lib in $${FOUND_LIBS}; do \
				name=`basename $${found_lib}`; \
				echo "# $${dir}: go get corelibs/$${name}"; \
				go get $${found_lib}@latest; \
			done; \
		fi; \
		popd > /dev/null; \
	done

tidy:
	@for dir in `find . -name "go.mod" -exec dirname \{\} \;`; do \
		pushd "$${dir}" > /dev/null; \
		echo "# $${dir}: go mod tidy"; \
		go mod tidy; \
		popd > /dev/null; \
	done

build:
	@for dir in `find . -name "go.mod" -exec dirname \{\} \;`; do \
		pushd "$${dir}" > /dev/null; \
		echo "# $${dir}: go build -v ./..."; \
		go build -v ./...; \
		popd > /dev/null; \
	done
