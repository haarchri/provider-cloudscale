
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
package_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

crossplane_bin = $(go_bin)/kubectl-crossplane

# Build kubectl-crossplane plugin
$(crossplane_bin):export GOBIN = $(go_bin)
$(crossplane_bin): | $(go_bin)
	go install github.com/crossplane/crossplane/cmd/crank@latest
	@mv $(go_bin)/crank $@

.PHONY: package
package: ## All-in-one packaging and releasing
package: package-push

.PHONY: package-provider
package-provider: export CONTROLLER_IMG = $(CONTAINER_IMG)
package-provider: $(crossplane_bin) generate-go ## Build Crossplane package
	@rm -rf package/*.xpkg
	@yq e '.spec.controller.image=strenv(CONTROLLER_IMG)' $(package_dir)/crossplane.yaml.template > $(package_dir)/crossplane.yaml
	@$(crossplane_bin) build provider -f $(package_dir)
	@echo Package file: $$(ls $(package_dir)/*.xpkg)

.PHONY: package-push
package-push: pkg_file = $(shell ls $(package_dir)/*.xpkg)
package-push: package-provider ## Push Crossplane package to container registry
	$(crossplane_bin) push provider -f $(pkg_file) $(PACKAGE_IMG)

.PHONY: .package-clean
.package-clean:
	rm -f $(crossplane_bin) package/*.xpkg $(package_dir)/crossplane.yaml
