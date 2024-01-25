FROM node:alpine as frontend-builder
WORKDIR /workdir
COPY ./frontend /workdir

RUN npm install -g pnpm && pnpm install && pnpm run build

FROM golang:1.21-alpine3.18 as builder

RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev
WORKDIR /workdir

RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=1 GOARCH=amd64 GOAMD64=v3 go build -trimpath -ldflags '-w -s' -o /server ./cmd/server

FROM alpine:latest

RUN apk add --no-cache libgcc

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /server /
COPY --from=frontend-builder /workdir/dist /frontend
COPY ./examples /frontend/assets
ENTRYPOINT ["/server"]
