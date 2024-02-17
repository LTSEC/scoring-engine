BUILD_NUMBER ?= dev+$(shell date -u '+%Y%m%d%H%M%S')
GO111MODULE = auto
export GO111MODULE

LDFLAGS = -w -s -X main.Build=$(BUILD_NUMBER)
BINARY_NAME=scoring-engine

ALL = linux-amd64

bin:
	templ generate ./web/*.templ
	go build -ldflags "$(LDFLAGS)"

build:
	mkdir -p bin/
	GOARCH=amd64 GOOS=windows go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-windows .
	GOARCH=amd64 GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux .
	GOARCH=386 GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-386 .
	GOARCH=arm GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-arm .
	GOARCH=arm64 GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-arm64 .
	GOARCH=riscv64 GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-riscv64 .
	GOARCH=ppc64 GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-ppc64 .
	GOARCH=ppc64le GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-ppc64le .
	GOARCH=s390x GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-s390x .
	GOARCH=mips GOOS=linux go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY_NAME)-linux-mips .

clean:
	rm -R ./bin/*

test:
	go test -v

docs:
	godoc -http=:6060

# test-cov-html:
# 	go test -coverprofile=coverage.out
# 	go tool cover -html=coverage.out

# bench:
# 	go test -bench=.

# bench-cpu:
# 	go test -bench=. -benchtime=5s -cpuprofile=cpu.pprof
# 	go tool pprof go-audit.test cpu.pprof

# bench-cpu-long:
# 	go test -bench=. -benchtime=60s -cpuprofile=cpu.pprof
# 	go tool pprof go-audit.test cpu.pprof

release: $(ALL:%=build/scoring-engine%.tar.gz)

build/%/scoring-engine: .FORCE
	GOOS=$(firstword $(subst -, , $*)) \
		GOARCH=$(word 2, $(subst -, ,$*)) \
		go build -trimpath -ldflags "$(LDFLAGS)" -o $@ .

build/scoring-engine%.tar.gz: build/%/scoring-engine
	tar -zcv -C build/$* -f $@ scoring-engine

.FORCE:
.PHONY: test test-cov-html bench bench-cpu bench-cpu-long bin release
.DEFAULT_GOAL := bin