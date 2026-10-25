# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS build
ARG TARGETOS=linux
ARG TARGETARCH
ARG VERSION=dev
WORKDIR /src
COPY . .
RUN mkdir -p /rootfs/usr/local/bin /rootfs/workspace /rootfs/tmp \
    && chmod 1777 /rootfs/tmp \
    && chown 65532:65532 /rootfs/workspace \
    && CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -trimpath -ldflags "-s -w -X github.com/jabahum/foundry/internal/cli.Version=${VERSION}" \
    -o /rootfs/usr/local/bin/foundry ./cmd/foundry

FROM scratch
ARG VERSION=dev
LABEL org.opencontainers.image.title="Foundry" \
      org.opencontainers.image.description="Offline project generation from five embedded foundations" \
      org.opencontainers.image.source="https://github.com/mohuganda/foundry" \
      org.opencontainers.image.version=$VERSION
COPY --from=build /rootfs/ /
USER 65532:65532
WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/foundry"]
CMD ["help"]
