.SHELL: /usr/bin/env bash -e

opa_flags = --server --config-file opa.yaml --log-level debug --log-format text

docker_flags = --rm --env-file .env --volume $(CURDIR)/opa/opa.yaml:/opa.yaml

opa-local:
	cd opa && opa run $(opa_flags)

opa-container:
	docker run $(docker_flags) docker.io/openpolicyagent/opa:0.70.0 \
		run $(opa_flags)

opa-container-custom:
	docker build -t acme/opa:custom opa/
	docker run $(docker_flags) acme/opa:custom run $(opa_flags)

gist-local:
	cd gist && go run .

gist-container:
	docker build -t gist gist/
	docker run $(docker_flags) gist
