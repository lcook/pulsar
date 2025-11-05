# SPDX-License-Identifier: BSD-2-Clause
#
# Copyright (c) Lewis Cook <lcook@FreeBSD.org>
.POSIX:

include Mk/config.mk

default: targets

build:
	@echo ">>> Building ${PROGRAM}@${VERSION} for ${OPSYS}"
	GOOS=${OPSYS:tl} ${GO_CMD} build ${GO_FLAGS} -o ${PROGRAM} cmd/pulsar/pulsar.go\
	     && strip -s ${PROGRAM}

clean:
	@echo ">>> Cleaning up project root directory"
	${GO_CMD} clean
	rm -f ${PROGRAM}

install: build
	@echo ">>> Installing ${PROGRAM} and configuration file"
.if !exists(${CFGDIR})
	@echo ">> Creating configuration directory"
	mkdir -p ${CFGDIR}
.endif
.if !exists(${CONFIG})
	@echo ">> No configuration file \`${CONFIG}\` found in project root directory"
	@echo "   You may use the example configuration \`config.yaml.example\` to get"
	@echo "   started.  Make sure to rename the example afterwards accordingly and"
	@echo "   reinstall, or copy to the directory \`${CFGDIR}\`"
	@sleep 4
.else
	install -m600 ${CONFIG} ${CFGDIR}
.endif
	install -m755 ${PROGRAM} ${SBINDIR}
 # Do not install the RC service script on
 # a non-FreeBSD host.
.if ${OPSYS:tl} == "freebsd"
	install -m755 ${RC} ${RCDIR}/${RC:C/\.in//}
.endif

container:
.if !exists(container/${OPSYS})
	@echo ">>> '${OPSYS}' is an unsupported operating system"
	@false
.endif
.if !exists(${PODMAN_CMD})
	@echo ">>> podman binary \`${PODMAN_CMD}\` not found on host"
	@echo "   Install the corresponding package and try again"
	@false
.endif
	@echo ">>> Building ${PROGRAM}@${VERSION} container image for ${OPSYS}"
	@${PODMAN_CMD} build ${PODMAN_ARGS} --file container/${OPSYS} --tag ${OCI_TAG} .

publish:
	@echo ">>> Publishing container image to ${OCI_TAG}"
	@${PODMAN_CMD} push ${OCI_TAG}

deinstall:
	@echo ">>> Deinstalling ${PROGRAM}"
	rm -rfv ${CFGDIR} ${SBINDIR}/${PROGRAM} ${RCDIR}/${RC:C/\.in//}

targets:
	@echo Targets: ${.ALLTARGETS:S/^default//:S/.END//:S/targets//}

update:
	@echo ">>> Updating and tidying up Go dependencies"
	${GO_CMD} get -u -v ./...
	${GO_CMD} mod tidy -v
	${GO_CMD} mod verify

test:
	@echo ">>> Running Go unit tests"
	${GO_CMD} test -v ./...

lint:
	@echo ">>> Linting Go files"
.if !exists(${GOLANGCI_CMD})
	@echo ">> golangci-lint binary \`${GOLANGCI_CMD}\` not found on host"
	@echo "   Check if \`PREFIX\` is set correctly and whether the accompanying package is installed"
	@false
.endif
	${GOLANGCI_CMD} run

format:
	@echo ">>> Formatting Go files"
	find . -name "*.go" -exec ${GOFMT_CMD} -w {} \;

.PHONY: build clean container default deinstall install registry publish targets update format lint test
