FROM golang:1.13 as builder
ARG VERSION
ARG GIT_COMMIT
WORKDIR /go/src/github.com/tommy-muehle/speakerdeck_exporter
COPY *.go ./
COPY vendor ./vendor
COPY go.mod go.sum ./
COPY collector ./collector
COPY client ./client
COPY http ./http
RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -ldflags "-X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT}" -o exporter .

FROM alpine:latest
RUN addgroup -S exporter && adduser -S -G exporter exporter
USER exporter
COPY --from=builder /go/src/github.com/tommy-muehle/speakerdeck_exporter/exporter /usr/bin/
ENTRYPOINT [ "/usr/bin/exporter" ]
HEALTHCHECK CMD [ "/usr/bin/exporter", "-version" ]
