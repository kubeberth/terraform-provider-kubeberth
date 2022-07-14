VERSION=0.11.0
OS=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: init
init:
	mkdir -p ~/.terraform.d/plugins

.PHONY: build
build:
	mkdir -p ~/.terraform.d/plugins/local/kubeberth/kubeberth/${VERSION}/${OS}_${ARCH}
	go build -o bin/terraform-provider-kubeberth_v${VERSION}
	cp bin/terraform-provider-kubeberth_v${VERSION} ~/.terraform.d/plugins/local/kubeberth/kubeberth/${VERSION}/${OS}_${ARCH}/
