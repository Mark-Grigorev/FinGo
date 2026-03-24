ARG GO_VERSION=1.26.1
ARG ALPINE_VERSION=3.21
ARG VERSION=dev
ARG PORT=8008
ARG PATH_TO_MAIN_FILE=cmd/*.go

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build

ARG VERSION
ARG PATH_TO_MAIN_FILE

WORKDIR /app

RUN apk add --no-cache gcc libc-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -mod=vendor -a installsuffix cgo -o app \
    -ldflags "-X 'main.version=${VERSION}'" ${PATH_TO_MAIN_FILE}

FROM alpine:${ALPINE_VERSION}
ARG PORT

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app

COPY --from=build --chown=appuser:appgroup /app/app /app/

USER appuser

EXPOSE ${PORT}

CMD ["/app/app"]