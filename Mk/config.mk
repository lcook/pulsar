# SPDX-License-Identifier: BSD-2-Clause
#
# Copyright (c) Lewis Cook <lcook@FreeBSD.org>
VERSION=	0.1.6
PROGRAM=	pulsar
RC=		${PROGRAM}.in
TOML=		config.toml

# Determine what operating-system we are using
# to build on and set an according localbase
# prefix.  We are not *too* concerned as to
# getting this right on a Linux-like host, as
# it primarily will be built/ran on FreeBSD.
#
# However, if and when necessary, can be
# amended from the commandline by setting
# the PREFIX variable e.g.,
#
# $ make PREFIX=/opt build
OPSYS!=		uname -s
.if ${OPSYS:tl} == "freebsd"
PREFIX?=	/usr/local
.elif ${OPSYS:tl} == "linux"
PREFIX?=	/usr
.endif

ETCDIR=		${PREFIX}/etc
BINDIR=		${PREFIX}/bin
SBINDIR=	${PREFIX}/sbin

RCDIR=		${ETCDIR}/rc.d
CFGDIR=		${ETCDIR}/${PROGRAM}

GO_CMD=		${BINDIR}/go
.if !exists(${GO_CMD})
.error ${.newline}>>> Go binary `${GO_CMD}` not found on host\
       ${.newline}    Check if `PREFIX` is set correctly and whether the accompanying Go package is installed
.endif
GOFMT_CMD=	${BINDIR}/gofmt
GOLANGCI_CMD=	${BINDIR}/golangci-lint
GIT_CMD=	${BINDIR}/git

GH_ACCOUNT?=	lcook
GH_PROJECT?=	${PROGRAM}

GO_MODULE=	github.com/${GH_ACCOUNT}/${GH_PROJECT}
GO_FLAGS=	-v -ldflags \
		"-s -w -X ${GO_MODULE}/internal/version.Build=${VERSION}"

.if exists(${.CURDIR}/.git) && exists(${GIT_CMD})
HASH!=		${GIT_CMD} rev-parse --short HEAD
BRANCH!=	${GIT_CMD} symbolic-ref HEAD | sed 's,refs/heads/,,'
DIRTY!=	${GIT_CMD} status --porcelain
.if ${DIRTY}
HASH:=		${HASH}-dirty
.endif
VERSION:=	${BRANCH}/${VERSION}-${HASH}
.endif

# Configuration for building podman containers.
PODMAN_CMD=	${BINDIR}/podman
OCI_REPO=	localhost
OCI_TAG=	${OCI_REPO}/${PROGRAM}:${HASH}
