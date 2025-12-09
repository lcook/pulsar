# SPDX-License-Identifier: BSD-2-Clause
#
# Copyright (c) Lewis Cook <lcook@FreeBSD.org>
.POSIX:
VERSION=	0.2.0
CONFIG=		config.yaml
GH_ACCOUNT=	lcook
GH_PROJECT=	pulsar

OSNAME=		${.MAKE.OS}
.if ${OSNAME} == FreeBSD
PREFIX?=	/usr/local
.elif ${OSNAME} == Linux
PREFIX?=	/usr
.else
.error ${.newline}=> ${OSNAME} is an unsupported OS
.endif

ETCDIR=		${PREFIX}/etc
CFGDIR=		${ETCDIR}/${GH_PROJECT}
.if ${OSNAME} == FreeBSD
RCDIR=		${ETCDIR}/rc.d
RCCFG=		${GH_PROJECT}.in
.endif
BINDIR=		${PREFIX}/bin
SBINDIR=	${PREFIX}/sbin

GIT_CMD=	${BINDIR}/git
.if exists(${.CURDIR}/.git)
.  if exists(${GIT_CMD})
GIT_REPO=
.  else
.error ${.newline}=> git directory found '${.CURDIR}/.git' but '${GIT_CMD}' ${.newline}   not found on the system. Check if `PREFIX` is set correctly and ${.newline}   whether the accompanying git package is installed
.  endif
.endif

GO_CMD=		${BINDIR}/go
GOLANGCI_CMD=	${BINDIR}/golangci-lint
PODMAN_CMD=	${BINDIR}/podman

BUILD_DEPENDS=		${GO_CMD}
CONTAINER_DEPENDS=	${PODMAN_CMD}

.if defined(GIT_REPO)
GIT_HASH!=	git rev-parse --short HEAD
GIT_BRANCH!=	git symbolic-ref HEAD 2>/dev/null | sed 's,refs/heads/,,'
GIT_DIRTY!=	git status --porcelain
.  if ${GIT_DIRTY}
GIT_HASH:=	${GIT_HASH}-dirty
.  endif
VERSION:=	${GIT_BRANCH}/${VERSION}-${GIT_HASH}
.endif

GO_MODULE=	github.com/${GH_ACCOUNT}/${GH_PROJECT}
GO_FLAGS=	-v -ldflags \
		"-s -w -X ${GO_MODULE}/internal/version.Build=${VERSION}"

.if ${OSNAME} == FreeBSD
PODMAN_ARGS=	--network=host
.endif
OCI_REPO?=	localhost
OCI_TAG=	${OCI_REPO}/${GH_PROJECT}:${GIT_HASH}
.if ${OCI_REPO} != localhost
OCI_TAG=	${OCI_REPO}/${GH_ACCOUNT}/${GH_PROJECT}:${GIT_HASH}
.endif

default: build

build: build-requirements
	@echo -------------------------------------------------------------------
	@echo ">>> Building ${GH_PROJECT}@${VERSION} for ${OSNAME}"
	@echo -------------------------------------------------------------------
	GOOS=${OSNAME:tl} ${GO_CMD} build ${GO_FLAGS} -o ${GH_PROJECT} cmd/pulsar/pulsar.go\
             && strip -s ${GH_PROJECT}
	@echo

run: build-requirements
	@echo -------------------------------------------------------------------
	@echo ">>> Running ${GH_PROJECT}@${VERSION}"
	@echo -------------------------------------------------------------------
	${GO_CMD} run cmd/pulsar/pulsar.go -V 2

clean:
	@echo -------------------------------------------------------------------
	@echo ">>> Cleaning up project root directory"
	@echo -------------------------------------------------------------------
	${GO_CMD} clean
	rm -f ${GH_PROJECT}
	@echo

install: build
	@echo -------------------------------------------------------------------
	@echo ">>> Installing ${GH_PROJECT}@${VERSION} and configuration file"
	@echo -------------------------------------------------------------------
.if !exists(${CFGDIR})
	mkdir -p ${CFGDIR}
.endif
.if exists(${CONFIG})
	@echo "=> No configuration file \`${CONFIG}\` found in project root directory"
	@echo "   You may use the example configuration \`config.yaml.example\` to get"
	@echo "   started.  Make sure to rename the example afterwards accordingly and"
	@echo "   reinstall, or copy to the directory \`${CFGDIR}\`"
	@sleep 4
.else
	install -m600 ${CONFIG} ${CFGDIR}
.endif
	install -m755 ${GH_PROJECT} ${SBINDIR}
.if ${OSNAME} == FreeBSD
	install -m755 ${RCCFG} ${RCDIR}/${RCCFG:C/\.in//}
.endif
	@echo

deinstall:
	@echo -------------------------------------------------------------------
	@echo ">>> Deinstalling ${GH_PROJECT}@${VERSION}"
	@echo -------------------------------------------------------------------
	rm -rfv ${CFGDIR} ${SBINDIR}/${GH_PROJECT} ${RCDIR}/${RCCFG:C/\.in//}
	@echo

container:
.if !exists(container/${OSNAME})
	@echo "=> '${OSNAME}' is an unsupported operating system"
	@false
.endif
	@echo -------------------------------------------------------------------
	@echo ">>> Building ${GH_PROJECT}@${VERSION} container image for ${OSNAME}"
	@echo -------------------------------------------------------------------
	${PODMAN_CMD} build ${PODMAN_ARGS} --file container/${OSNAME} --tag ${OCI_TAG} .
	@echo

container-publish: container-requirements
	@echo -------------------------------------------------------------------
	@echo ">>> Publishing container image to ${OCI_TAG}"
	@echo -------------------------------------------------------------------
	@${PODMAN_CMD} push ${OCI_TAG}
	@echo

update:
	@echo -------------------------------------------------------------------
	@echo ">>> Updating and tidying up Go dependencies"
	@echo -------------------------------------------------------------------
	${GO_CMD} get -u -v ./...
	${GO_CMD} mod tidy -v
	${GO_CMD} mod verify
	@echo

test:
	@echo -------------------------------------------------------------------
	@echo ">>> Running Go unit tests"
	@echo -------------------------------------------------------------------
	${GO_CMD} test -v ./...
	@echo

lint:
	@echo -------------------------------------------------------------------
	@echo ">>> Linting Go files"
	@echo -------------------------------------------------------------------
.if !exists(${GOLANGCI_CMD})
	@echo "=> golangci-lint binary \`${GOLANGCI_CMD}\` not found on host"
	@echo "   Check if \`PREFIX\` is set correctly and whether the accompanying package is installed"
	@false
.endif
	${GOLANGCI_CMD} run
	@echo

build-requirements:
.for dep in ${BUILD_DEPENDS}
.  if !exists(${dep})
	@echo "=> Build dependency '${dep}' not found. Check if `PREFIX` is"
	@echo "   set correctly and whether the accompanying package is installed"
	@false
.  endif
.endfor

container-requirements:
.for dep in ${CONTAINER_DEPENDS}
.  if !exists(${dep})
	@echo "=> Container dependency '${dep}' not found. Check if \`PREFIX\` is"
	@echo "   set correctly and whether the accompanying package is installed"
	@false
.  endif
.endfor

.PHONY:	build run clean install deinstall container container-publish update test lint build-requirements container-requirements
