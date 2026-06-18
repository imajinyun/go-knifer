.PHONY: help doctor install-hooks uninstall-hooks test test-race coverage-profile coverage-report coverage-check api-check mod-verify tidy-check diff-whitespace diff-clean diff-check vet arch lint govulncheck quick-check security-check full-check bench bench-core bench-facade bench-codec bench-smoke check ci-test

GO ?= go
GOLANGCI_LINT ?= golangci-lint
PKGS ?= ./...
COVERAGE_FILE ?= coverage.out
BENCH ?= .
BENCH_PKGS ?= ./internal/slice ./internal/maps ./internal/str ./internal/num
BENCH_FACADE_PKGS ?= ./vslice ./vmap ./vstr ./vnum
BENCH_CODEC_PKGS ?= ./internal/json ./vjson ./internal/xml ./vxml
BENCHTIME ?= 1s
BENCHCOUNT ?= 1

help:
	@echo "Targets:"
	@echo "  test            Run unit tests"
	@echo "  test-race       Run race/shuffle tests and write coverage"
	@echo "  coverage-profile Generate race/shuffle coverage profile"
	@echo "  coverage-report  Print function coverage from COVERAGE_FILE"
	@echo "  coverage-check  Enforce repository and package coverage gates"
	@echo "  bench-core      Run core benchmark baselines"
	@echo "  bench-facade    Run facade benchmark baselines"
	@echo "  bench-codec     Run JSON/XML benchmark baselines"
	@echo "  bench-smoke     Run a short core benchmark smoke check"
	@echo "  doctor          Diagnose local Go/tooling/Git environment"
	@echo "  install-hooks   Enable optional local Git validation hooks"
	@echo "  uninstall-hooks Disable optional local Git validation hooks"
	@echo "  quick-check     Run fast local governance gates"
	@echo "  security-check  Run lint and govulncheck"
	@echo "  full-check      Run full local gates with race coverage"
	@echo "  api-check       Verify exported API snapshot is current"
	@echo "  check           Run local stability gates"
	@echo "  ci-test         Run CI test-job gates"

doctor:
	@echo "== go =="
	@$(GO) version
	@echo "== python3 =="
	@if command -v python3 >/dev/null 2>&1; then python3 --version; else echo "python3 not found"; fi
	@echo "== golangci-lint =="
	@if command -v $(GOLANGCI_LINT) >/dev/null 2>&1; then $(GOLANGCI_LINT) version; else echo "$(GOLANGCI_LINT) not found"; fi
	@echo "== govulncheck =="
	@if $(GO) tool govulncheck -version >/dev/null 2>&1; then $(GO) tool govulncheck -version; else echo "govulncheck tool not available; run: go get -tool golang.org/x/vuln/cmd/govulncheck"; fi
	@echo "== git status =="
	@git status --short --branch
	@echo "== module =="
	@$(GO) list -m | grep 'go-knifer' || $(GO) list -m
	@echo "== package list =="
	@$(GO) list ./... >/dev/null
	@echo "go list ./... OK"

install-hooks:
	@git config core.hooksPath .githooks
	@chmod +x .githooks/pre-commit .githooks/pre-push
	@echo "Installed go-knifer Git hooks from .githooks"

uninstall-hooks:
	@git config --unset core.hooksPath || true
	@echo "Disabled go-knifer Git hooks"

test:
	$(GO) test $(PKGS)

test-race:
	$(GO) test -race -shuffle=on -coverprofile=$(COVERAGE_FILE) $(PKGS)

coverage-profile: test-race

coverage-report:
	$(GO) tool cover -func=$(COVERAGE_FILE)

coverage-check:
	bash bin/check_coverage.sh $(COVERAGE_FILE)

api-check:
	bash bin/check_api_compat.sh

mod-verify:
	$(GO) mod verify

tidy-check:
	$(GO) mod tidy
	git diff --exit-code -- go.mod go.sum

diff-whitespace:
	git diff --check

diff-clean:
	git diff --exit-code

diff-check: diff-whitespace diff-clean

vet:
	$(GO) vet $(PKGS)

arch:
	bash bin/check_arch.sh

lint:
	$(GOLANGCI_LINT) run $(PKGS)

govulncheck:
	$(GO) tool govulncheck $(PKGS)

quick-check: mod-verify vet arch test api-check diff-whitespace

security-check: lint govulncheck

full-check: mod-verify vet arch test-race coverage-check api-check lint govulncheck diff-whitespace

bench:
	$(GO) test -bench=$(BENCH) -benchmem -benchtime=$(BENCHTIME) -count=$(BENCHCOUNT) -run=^$$ $(BENCH_PKGS)

bench-core: bench

bench-facade:
	$(GO) test -bench=$(BENCH) -benchmem -benchtime=$(BENCHTIME) -count=$(BENCHCOUNT) -run=^$$ $(BENCH_FACADE_PKGS)

bench-codec:
	$(GO) test -bench=$(BENCH) -benchmem -benchtime=$(BENCHTIME) -count=$(BENCHCOUNT) -run=^$$ $(BENCH_CODEC_PKGS)

bench-smoke:
	$(GO) test -bench=Benchmark -benchtime=100ms -count=1 -run=^$$ $(BENCH_PKGS) $(BENCH_FACADE_PKGS) $(BENCH_CODEC_PKGS)

check: full-check

ci-test: mod-verify vet tidy-check diff-check arch test-race coverage-check api-check
