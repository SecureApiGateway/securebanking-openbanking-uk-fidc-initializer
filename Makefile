service := secureopenbanking-uk-iam-initializer
gcr-repo := sbat-gcr-develop
binary-name := initialize
github-tag := v1.0.0
config-github-sha := e99bfc44ff4d66e70eef1be61363887d383a8aa7

define clone_config_tag
	if [ -d "config-repo" ]; then rm config-repo -rf; fi
	if [ -d "config" ]; then rm config -rf; fi
    git clone --depth 1 --branch ${github-tag} git@github.com:SecureApiGateway/fr-platform-config.git config-repo
	mkdir config
	mv config-repo/config/* config
	rm -rf config-repo
endef

define clone_config
	if [ -d "config-repo" ]; then rm config-repo -rf; fi
	if [ -d "config" ]; then rm config -rf; fi
    git clone git@github.com:SecureApiGateway/fr-platform-config.git config-repo
	cd config-repo && git checkout ${config-github-sha}
	mkdir config
	mv config-repo/config/* config
	rm -rf config-repo
endef

.PHONY: all
all: mod build

mod:
	go mod download

build: clean mod
	$(call clone_config)
	go build -o ${binary-name}

test:
	go test ./...

test-ci: mod
	$(eval localPath=$(shell pwd))
	curl -fsSL https://raw.githubusercontent.com/pact-foundation/pact-ruby-standalone/master/install.sh | bash
	PATH=$(PATH):${localPath}/pact/bin go test ./...

clean:
	go clean
	rm -f ${binary-name}

docker: clean mod
ifndef tag
	$(warning No tag supplied, latest assumed. supply tag with make docker tag=x.x.x service=...)
	$(eval tag=latest)
endif
	env GOOS=linux GOARCH=amd64 go build -o initialize
	docker buildx build --platform linux/amd64  -t eu.gcr.io/${gcr-repo}/securebanking/${service}:${tag} .
	docker push eu.gcr.io/${gcr-repo}/securebanking/${service}:${tag}
ifdef release-repo
	docker tag eu.gcr.io/${gcr-repo}/securebanking/${service}:${tag} eu.gcr.io/${release-repo}/securebanking/${service}:${tag}
	docker push eu.gcr.io/${release-repo}/securebanking/${service}:${tag}
endif
