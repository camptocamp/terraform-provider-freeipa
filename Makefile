TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
WEBSITE_REPO=github.com/hashicorp/terraform-website
PKG_NAME=freeipa

# Local dev env creds
export FREEIPA_HOST=ipa.example.test
export FREEIPA_USERNAME=admin
export FREEIPA_PASSWORD=SecretPassword123
export FREEIPA_INSECURE=true

default: build

build: fmtcheck
	CGO_ENABLED=0 go install

prepare-dev-env:
	@command -v podman >/dev/null 2>&1 || { echo >&2 "ERROR: podman is not installed. Please install it before running this make target."; exit 1; }
	@if ! grep -qx '127\.0\.0\.1\sipa.example.test' /etc/hosts; then \
		echo >&2 "WARNING: The hostname 'ipa.example.test' is not correctly redirected to 127.0.0.1 in your /etc/hosts file."; \
		echo >&2 "         Please add the following line to your /etc/hosts file and try again:"; \
		echo >&2 "         127.0.0.1 ipa.example.test"; \
		exit 1; \
	fi;
	@if ! podman ps --format '{{.Names}}' | grep -w freeipa-server &> /dev/null; then \
		echo "Starting the freeipa container environment..."; \
		if ! podman compose up -d; then \
			echo >&2 "ERROR: Failed to start the freeipa container environment."; \
			exit 1; \
		fi; \
		echo "Waiting for the freeipa container environment to come up..."; \
		until curl -s -k https://ipa.example.test/ipa/ui/ | grep -q "Identity Management"; do sleep 1; done; \
		echo "FreeIPA container environment started."; \
	else \
		echo "FreeIPA container environment already running."; \
	fi

stop-dev-env:
	podman compose down

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"


test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build prepare-dev-env stop-dev-env test testacc vet fmt fmtcheck errcheck test-compile website website-test
