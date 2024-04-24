service := secureopenbanking-uk-iam-initializer
repo := europe-west4-docker.pkg.dev/sbat-gcr-develop/sapig-docker-artifact
binary-name := initialize


.PHONY: all
all: mod build
	
mod:
	go mod download

build: clean mod
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
	$(warning no tag supplied; latest assumed)
	$(eval tag=latest)
endif
ifndef setlatest
	$(warning no setlatest true|false supplied; false assumed)
	$(eval setlatest=false)
endif
	env GOOS=linux GOARCH=amd64 go build -o initialize
	if [ "${setlatest}" = "true" ]; then \
		docker buildx build --platform linux/amd64 -t ${repo}/securebanking/${service}:${tag} -t ${repo}/securebanking/${service}:latest . ; \
		docker push ${repo}/securebanking/${service} --all-tags; \
    else \
   		docker buildx build --platform linux/amd64 -t ${repo}/securebanking/${service}:${tag} . ; \
   		docker push ${repo}/securebanking/${service}:${tag}; \
   	fi;

ifdef release-repo
	docker tag ${repo}/securebanking/${service}:${tag} ${release-repo}/securebanking/${service}:${tag}
	docker push ${release-repo}/securebanking/${service}:${tag}
endif
