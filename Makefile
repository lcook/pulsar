# SPDX-License-Identifier: BSD-2-Clause
#
# Copyright (c) Lewis Cook <lcook@FreeBSD.org>
.POSIX:

include Mk/config.mk

default: targets

build:
	@echo ">>> Building ${PROGRAM}@${VERSION} for ${OPSYS}"
	GOOS=${OPSYS:tl} ${GO_CMD} build ${GO_FLAGS} -o ${PROGRAM}\
	     && strip -s ${PROGRAM}

clean:
	@echo ">>> Cleaning up project root directory"
	${GO_CMD} clean

install: build
	@echo ">>> Installing ${PROGRAM} and configuration file"
.if !exists(${CFGDIR})
	@echo ">> Creating configuration directory"
	mkdir -p ${CFGDIR}
.endif
.if !exists(${TOML})
	@echo ">> No configuration file \`${TOML}\` found in project root directory"
	@echo "   You may use the example configuration \`config.toml.example\` to get"
	@echo "   started.  Make sure to rename the example afterwards accordingly and"
	@echo "   reinstall, or copy to the directory \`${CFGDIR}\`"
	@sleep 4
.else
	install -m600 ${TOML} ${CFGDIR}
.endif
	install -m755 ${PROGRAM} ${SBINDIR}
        # Do not install the RC service script on
        # a non-FreeBSD host.
.if ${OPSYS:tl} == "freebsd"
	install -m755 ${RC} ${RCDIR}/${RC:C/\.in//}
.endif

deinstall:
	@echo ">>> Deinstalling ${PROGRAM}"
	rm -rfv ${CFGDIR} ${SBINDIR}/${PROGRAM} ${RCDIR}/${RC:C/\.in//}

targets:
	@echo Targets: ${.ALLTARGETS:S/^default//:S/.END//:S/targets//}

update:
	@echo ">>> Updating and tidying up Go dependencies"
	${GO_CMD} get -u -v
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

.PHONY: build clean default deinstall install targets update format lint test
