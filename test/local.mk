setup_envtest_bin = $(kind_dir)/setup-envtest

# Prepare binary
# We need to set the Go arch since the binary is meant for the user's OS.
$(setup_envtest_bin): export GOOS = $(shell go env GOOS)
$(setup_envtest_bin): export GOARCH = $(shell go env GOARCH)
$(setup_envtest_bin):
	@mkdir -p $(kind_dir)
	cd test && go build -o $@ sigs.k8s.io/controller-runtime/tools/setup-envtest
	$@ $(ENVTEST_ADDITIONAL_FLAGS) use '$(ENVTEST_K8S_VERSION)!'
	chmod -R +w $(kind_dir)/k8s

ifeq ($(shell uname -s),Darwin)
	b64 := base64
else
	b64 := base64 -w0
endif

.PHONY: local-install
local-install: export KUBECONFIG = $(KIND_KUBECONFIG)
local-install: kind-load-image install-crd ## Install Operator in local cluster
	helm upgrade --install appcat-service-s3 charts/appcat-service-s3 \
		--create-namespace --namespace appcat-service-s3-system \
		--set "operator.args[0]=--log-level=2" \
		--set "operator.args[1]=operator" \
		--set podAnnotations.date="$(shell date)" \
		--wait $(local_install_args)

.PHONY: kind-run-operator
kind-run-operator: export KUBECONFIG = $(KIND_KUBECONFIG)
kind-run-operator: kind-setup ## Run in Operator mode against kind cluster (you may also need `install-crd`)
	go run . -v 1 operator

###
### Integration Tests
###

.PHONY: test-integration
test-integration: export ENVTEST_CRD_DIR = $(shell realpath $(envtest_crd_dir))
test-integration: $(setup_envtest_bin) .envtest_crds ## Run integration tests against code
	export KUBEBUILDER_ASSETS="$$($(setup_envtest_bin) $(ENVTEST_ADDITIONAL_FLAGS) use -i -p path '$(ENVTEST_K8S_VERSION)!')" && \
	go test -tags=integration -coverprofile cover.out -covermode atomic ./...

envtest_crd_dir ?= $(kind_dir)/crds

.envtest_crd_dir:
	@mkdir -p $(envtest_crd_dir)
	@cp -r package/crds $(kind_dir)

.envtest_crds: .envtest_crd_dir