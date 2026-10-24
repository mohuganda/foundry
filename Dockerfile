# syntax=docker/dockerfile:1
FROM --platform=$BUILDPLATFORM golang:1.23-alpine AS build
ARG TARGETOS=linux
ARG TARGETARCH
ARG VERSION=dev
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -trimpath -ldflags "-s -w -X github.com/jabahum/foundry/internal/cli.Version=${VERSION}" \
    -o /out/foundry ./cmd/foundry \
    && mkdir -p /out/workspace

FROM scratch
ARG VERSION=dev
LABEL org.opencontainers.image.title="Foundry" \
      org.opencontainers.image.description="Offline project generation from five embedded foundations" \
      org.opencontainers.image.source="https://github.com/mohuganda/foundry" \
      org.opencontainers.image.version=$VERSION
COPY --from=build /out/foundry /usr/local/bin/foundry
COPY --from=build --chown=65532:65532 /out/workspace /workspace
USER 65532:65532
WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/foundry"]
CMD ["help"]
