# Copyright (c) 2021, Lewis Cook <lcook@FreeBSD.org>
#
# Targets intended for use on the command line
#
# default	- Runs `build` target
# build         - Build project 
# install	- Install `pulseline` and configuration globally
# deinstall	- Remove all files installed from `install` target
# clean		- Cleanup any unnecessary files from building the project
# target        - Print all available targets
#
# Targets intended for managing go related tasks
#
# format	- Format Go files with `gofmt`
# lint		- Run `golangci-lint` across project source
# test		- Run go unit tests
# mod		- Download required Go modules needed to build
# mod-update    - Updates Go modules
VERSION=	0.1.4
PROGRAM=	pulseline
RC=		${PROGRAM}.in
YAML=		config.yaml

LOCALBASE?=	/usr/local

ETCDIR=		${LOCALBASE}/etc
BINDIR=		${LOCALBASE}/bin
SBINDIR=	${LOCALBASE}/sbin
RCDIR=		${ETCDIR}/rc.d
CFGDIR=		${ETCDIR}/${PROGRAM}

GO_CMD=		${BINDIR}/go
GOFMT_CMD=	${BINDIR}/gofmt
GOLANGCI_CMD=	${BINDIR}/golangci-lint
GIT_CMD=	${BINDIR}/git

GO_FLAGS=	-v -ldflags "-s -w -X main.Version=${VERSION}"

.if !exists(${GO_CMD})
.error ${.newline}WARNING: go not installed. Install by running ${.newline}pkg install lang/go.
.endif

.if exists(${.CURDIR}/.git) && exists(${GIT_CMD})
HASH!=		${GIT_CMD} rev-parse --short HEAD
BRANCH!=	${GIT_CMD} symbolic-ref HEAD | sed 's,refs/heads/,,'
VERSION:=	${BRANCH}/${VERSION}-${HASH}
.endif

default: build .PHONY

build: .PHONY
	@echo "--------------------------------------------------------------"
	@echo ">>> Building ${PROGRAM}@${VERSION}"
	@echo "--------------------------------------------------------------"
	${GO_CMD} build ${GO_FLAGS} -o ${PROGRAM} && \
		strip -s ${PROGRAM}

install: build .PHONY
	@echo "--------------------------------------------------------------"
	@echo ">>> Installing ${PROGRAM} and configuration"
	@echo "--------------------------------------------------------------"
.if !exists(${CFGDIR})
	mkdir -p ${CFGDIR}
.endif
.if !exists(${YAML})
	@echo "WARNING: Configuration file (${YAML}) not found in"
	@echo "projecct root directory.  Use the example configuration"
	@echo "(config.example.yaml) to get started AND copy to"
	@echo "${CFGDIR}."
	@sleep 3
.else
	install -m600 ${YAML} ${CFGDIR}
.endif
	install -m755 ${PROGRAM} ${SBINDIR}
	install -m755 ${RC} ${RCDIR}/${RC:C/\.in//}

deinstall: .PHONY
	@echo "--------------------------------------------------------------"
	@echo ">>> Deinstalling ${PROGRAM}"
	@echo "--------------------------------------------------------------"
	rm -rfv ${CFGDIR} ${SBINDIR}/${PROGRAM} ${RCDIR}/${RC:C/\.in//}

clean: .PHONY
	@echo "--------------------------------------------------------------"
	@echo ">>> Cleaning up project root directory"
	@echo "--------------------------------------------------------------"
	${GO_CMD} clean

targets help: .PHONY
	@echo Targets: ${.ALLTARGETS:S/^default//:S/.END//}

mod: .PHONY
	${GO_CMD} mod tidy -v
	${GO_CMD} mod verify

mod-update: .PHONY
	${GO_CMD} get -u -v

test: .PHONY
	@echo "--------------------------------------------------------------"
	@echo ">>> Running unit tests"
	@echo "--------------------------------------------------------------"
	${GO_CMD} test -v ./...

lint: .PHONY
.if !exists(${GOLANGCI_CMD})
	@echo "WARNING: golangci-lint not installed.  Install by running"
	@echo "pkg install devel/golangci-lint."
	@sleep 3
	@false
.endif
	${GOLANGCI_CMD} run

format: .PHONY
	find . -name "*.go" -exec ${GOFMT_CMD} -w {} \;
