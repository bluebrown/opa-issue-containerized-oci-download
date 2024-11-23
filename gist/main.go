package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/containerd/errdefs"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	oraslib "oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/oci"
)

func main() {
	tlsConfig := tls.Config{}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.ResponseHeaderTimeout = time.Duration(30) * time.Second
	transport.TLSClientConfig = &tlsConfig

	httpClient := *http.DefaultClient
	httpClient.Transport = transport

	urlInfo, err := url.Parse(os.Getenv("ACR_URL"))
	if err != nil {
		log.Fatalf("parse url: %v", err)
	}

	acrToken := base64.StdEncoding.EncodeToString([]byte(os.Getenv("ACR_TOKEN")))

	authorizer := docker.NewDockerAuthorizer(
		docker.WithAuthHeader(http.Header{"Authorization": {"Basic " + acrToken}}),
		docker.WithAuthClient(&httpClient),
	)

	registryHost := docker.RegistryHost{
		Host:         urlInfo.Host,
		Scheme:       urlInfo.Scheme,
		Capabilities: docker.HostCapabilityPull | docker.HostCapabilityResolve | docker.HostCapabilityPush,
		Client:       &httpClient,
		Path:         "/v2",
		Authorizer:   authorizer,
	}

	opts := docker.ResolverOptions{
		Hosts: func(string) ([]docker.RegistryHost, error) {
			return []docker.RegistryHost{registryHost}, nil
		},
	}

	resolver := docker.NewResolver(opts)

	manager := remoteManager{
		resolver: resolver,
		srcRef:   os.Getenv("ACR_BUNDLE"),
	}

	localstore, err := oci.New(filepath.Join(os.TempDir(), "opa", "oci"))
	if err != nil {
		log.Fatalf("oci store: %v", err)
	}

	spec, err := oraslib.Copy(context.Background(), &manager,
		manager.srcRef, localstore, "", oraslib.DefaultCopyOptions)
	if err != nil {
		log.Fatalf("oras copy: %v", err)
	}

	log.Printf("ok: copied %d bytes: %s", spec.Size, spec.MediaType)
}

type remoteManager struct {
	resolver remotes.Resolver
	srcRef   string
}

func (r *remoteManager) Resolve(ctx context.Context, ref string) (ocispec.Descriptor, error) {
	_, desc, err := r.resolver.Resolve(ctx, ref)
	if err != nil {
		return ocispec.Descriptor{}, err
	}
	return desc, nil
}

func (r *remoteManager) Fetch(ctx context.Context, target ocispec.Descriptor) (io.ReadCloser, error) {
	fetcher, err := r.resolver.Fetcher(ctx, r.srcRef)
	if err != nil {
		return nil, err
	}
	return fetcher.Fetch(ctx, target)
}

func (r *remoteManager) Exists(ctx context.Context, target ocispec.Descriptor) (bool, error) {
	_, err := r.Fetch(ctx, target)
	if err == nil {
		return true, nil
	}

	return !errdefs.IsNotFound(err), err
}
