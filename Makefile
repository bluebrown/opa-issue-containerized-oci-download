opa_flags ?= --server --config-file /opa.yaml --log-level debug --log-format text
docker_flags ?= --rm --env-file .env --volume $(CURDIR)/opa/opa.yaml:/opa.yaml

opa-local:
	opa run -c opa/opa.yaml

opa-container:
	docker run $(docker_flags) docker.io/openpolicyagent/opa:0.70.0 run $(opa_flags)

opa-container-custom:
	docker build -t acme/opa:custom opa/
	docker run $(docker_flags) acme/opa:custom run $(opa_flags)

gist-local:
	cd gist && go run .

gist-container:
	docker run $(docker_flags) docker build -t gist gist/
	docker run --rm --env-file .env gist
