FROM golang:1.22.5-alpine AS builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash git make gcc gettext musl-dev

#dependencies
COPY ["go.mod", "go.sum", "./"]

RUN go mod download

#copy
COPY cmd ./cmd
COPY config ./config
COPY internal ./internal

#build
RUN go build -o ./bin/app ./cmd/url-shortener/main.go

FROM alpine AS runner

COPY --from=builder /usr/local/src/bin/app /
COPY --from=builder /usr/local/src/internal/storage/storage.db /internal/storage/storage.db
COPY --from=builder /usr/local/src/config/local.yaml /local.yaml

CMD ["/app"]