FROM docker.io/library/debian
RUN apt-get update && apt-get install -y ca-certificates
ADD https://github.com/open-policy-agent/opa/releases/download/v0.70.0/opa_linux_amd64 .
RUN echo '2879c01f1e5762f28e27c9f81b4035bd5f532753f18c2c6dcbc2943347cc6ea5 opa_linux_amd64' | sha256sum -c
RUN install opa_linux_amd64 /usr/local/bin/opa
ENTRYPOINT ["/usr/local/bin/opa"]

