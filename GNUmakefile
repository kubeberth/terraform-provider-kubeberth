VERSION=0.0.5

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
	mkdir -p ~/.terraform.d/plugins/local/kubeberth/kubeberth/${VERSION}/linux_amd64
	go build -o bin/terraform-provider-kubeberth_v${VERSION}
	cp bin/terraform-provider-kubeberth_v${VERSION} ~/.terraform.d/plugins/local/kubeberth/kubeberth/${VERSION}/linux_amd64
