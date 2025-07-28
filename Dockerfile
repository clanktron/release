FROM golang:1.24-alpine AS builder-base
RUN apk add --no-cache ca-certificates make git
WORKDIR /build
COPY go.mod /build/
RUN --mount=type=cache,target=/go/pkg/mod go mod download

FROM builder-base AS builder
COPY . /build
ARG BINOS
ARG BINARCH
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    make GOOS=${BINOS} GOARCH=${BINARCH}

FROM scratch AS bin
ARG BINOS
ARG BINARCH
COPY --from=builder /build/git-release /git-release-${BINOS}-${BINARCH}

FROM scratch AS production
COPY <<EOF /etc/passwd
"git-release:x:1001:1001::/:"
EOF
COPY <<EOF /etc/group
"git-release:x:2000:git-release"
EOF
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/git-release /
USER git-release
ENTRYPOINT ["/git-release"]
