GOPROXY := https://goproxy.cn,direct
GHPROXY := ""

.PHONY: build clean
build:
	./build.sh build
	@echo "" && make format

clean:
	@./build.sh clean

format:
	@echo "run gofmt && goimports"
	@GOFILES=$$(find . -name "*.go"); \
	for path in $${GOFILES}; do \
		gofmt -w -l $${path}; \
		goimports -w -l $${path}; \
	done;

auto-update-infra-for-master:
	@if [[ "$$(git rev-parse --abbrev-ref HEAD)" == "master" ]]; then \
		make update-infra; \
	fi;

update-infra:
	echo "update infra and gohub"
	go env -w GO111MODULE=on
	go env -w GOPROXY="$(GOPROXY)"
	go get -u github.com/erda-project/erda-infra
	go get -u github.com/erda-project/erda-infra/tools/gohub
	gohub tools install -v --goproxy="$(GOPROXY)" --ghproxy="$(GHPROXY)"

build-use-docker-image:
	proj_root=$$(cd ../../ && pwd) && \
	docker run --rm -ti \
		-v $${proj_root}:/go/src/github.com/erda-project/erda \
		-v $$(go env GOMODCACHE):/go/pkg/mod \
		-e target_proj_root=$${target_proj_root} \
		registry.erda.cloud/erda/gohub:latest \
		sh -c 'cd /go/src/github.com/erda-project/erda && make proto-go-in-local'
