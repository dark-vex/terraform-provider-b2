NAME=py-terraform-provider-b2
DIR=b2_terraform
EGG_INFO=${NAME}.egg-info
SPEC=${NAME}.spec
OS=$(shell python -c "import platform; print(platform.system())")
STATICX?=1

default: build

.PHONY: deps deps-test format lint test testacc clean build

deps:
	@pip install -r requirements-dev.txt

deps-test:
	@pip install -r requirements.txt

format:
	@ruff format ${DIR} tests

lint:
	@ruff check ${DIR} tests
	@python ../scripts/check-headers.py '**/*.py'

test:
	@PYTHONPATH=. python -m unittest discover -s tests -v

testacc: build

clean:
	@rm -rf build dist ${EGG_INFO}

build:
	@pyinstaller ${SPEC}
ifeq ($(OS)$(STATICX), Linux1)
	@mv -f dist/py-terraform-provider-b2 dist/py-terraform-provider-b2-linked
	@staticx --strip --loglevel INFO dist/py-terraform-provider-b2-linked dist/py-terraform-provider-b2
	@rm -f dist/py-terraform-provider-b2-linked
endif
ifeq ($(OS), Windows)
	@mv -f dist/py-terraform-provider-b2.exe dist/py-terraform-provider-b2
endif
ifneq ($(origin CI), undefined)
	@echo "version=$(subst refs/tags/v,,${GITHUB_REF})" > "${GITHUB_OUTPUT}"
endif

install: build

docs: build

docs-lint: build

all: deps lint build
