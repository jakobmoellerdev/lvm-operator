# https://docs.docker.com/engine/reference/builder/#automatic-platform-args-in-the-global-scope
ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM
# Build the manager binary
FROM golang:1.21 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# since we use vendoring we don't need to redownload our dependencies every time. Instead we can simply
# reuse our vendored directory and verify everything is good. If not we can abort here and ask for a revendor.
COPY vendor vendor/
RUN go mod verify

# Copy the go source
COPY api/ api/
COPY cmd/ cmd/
COPY internal/ internal/

# Build
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -mod=vendor --ldflags "-s -w" -a -o lvms cmd/main.go

FROM --platform=$TARGETPLATFORM registry.ci.openshift.org/ocp/4.15:base-rhel9 as baseocp
# vgmanager needs 'nsenter' and other basic linux utils to correctly function
FROM --platform=$TARGETPLATFORM registry.access.redhat.com/ubi9/ubi-minimal:9.2

COPY --from=baseocp /etc/yum.repos.d/localdev-rhel-9-baseos-rpms.repo /etc/yum.repos.d/localdev-rhel-9-baseos-rpms.repo
COPY --from=baseocp /etc/yum.repos.d/redhat.repo /etc/yum.repos.d/redhat.repo
RUN microdnf update -y && microdnf install -y util-linux e2fsprogs xfsprogs glibc && microdnf clean all

WORKDIR /
COPY --from=builder /workspace/lvms .
USER 65532:65532

# '/lvms' is the entrypoint for all LVMS binaries
ENTRYPOINT ["/lvms"]
