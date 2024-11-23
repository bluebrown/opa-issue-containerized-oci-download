# OPA Issue - OCI download in container image

## OPA

OPA fails to fetch OCI resources when run in docker. The [opa
directory](./opa/) contains the files needed to reproduce the error.

When OPA runs locally, it does not fail.

```console
$ opa version
Version: 0.70.0
Build Commit: 2ea031ea04e6a8afbc5dd22f656131dc3cfc5a7d
Build Timestamp: 2024-10-31T19:39:52Z
Build Hostname: 799a5774bce7
Go Version: go1.23.1
Platform: linux/amd64
WebAssembly: available

$ opa run -c opa/opa.yaml -s --log-format text
[INFO] Bundle loaded and activated successfully. Etag updated to 2c2def83472ef874c864e9a9af91e57cffef3076537e34919cfb5b2d459d6b7d.
  name = "main"
  plugin = "bundle"
```

When running OPA in docker, it fails. Both with the offical OPA image
from dockerhub and the custom image build from the
[Dockerfile](./opa/Dockerfile) in the opa directory. Below outout is
from the official image.

```console
$ alias dopa='docker run --rm --env-file .env \
    -v "$PWD/opa/opa.yaml:/opa.yaml" \
    docker.io/openpolicyagent/opa:0.70.0'

$ dopa version
Version: 0.70.0
Build Commit: 2ea031ea04e6a8afbc5dd22f656131dc3cfc5a7d
Build Timestamp: 2024-10-31T19:39:52Z
Build Hostname: fc72eb07593c
Go Version: go1.23.1
Platform: linux/amd64
WebAssembly: available

$ dopa run -c /opa.yaml -s --log-format text
[ERROR] Bundle load failed: failed to pull
org.azurecr.io/acme/policy:latest: download for
'org.azurecr.io/acme/policy:latest' failed: failed to
ingest: copy failed: httpReadSeeker: failed open: unexpected status code
https://org.azurecr.io/v2/acme/policy/blobs/sha256:ca3d163bab055381827226140568f3bef7eaac187cebd76878e0b63e9e442356:
403 Server failed to authenticate the request. Make sure the value of
Authorization header is formed correctly including the signature.
  name = "main"
  plugin = "bundle"
```

## Gist

The [gist directory](./gist/), contains go source code, that replicates
opas rest authentication and oci downloader logic. It uses the same
libraries and tried to configure the components in the same fashion as
opa.

```console
$ cd gist
$ go run .
2024/11/23 16:21:20 ok: copied 549 bytes: application/vnd.oci.image.manifest.v1+json
```

It runs succesfully in a container image.

```console
$ docker build -t gist gist/
$ docker run --rm --env-file .env gist
2024/11/23 16:48:19 ok: copied 549 bytes: application/vnd.oci.image.manifest.v1+json
```
