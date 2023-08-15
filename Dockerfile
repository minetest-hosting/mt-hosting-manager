FROM golang:1.20.7 as go-app
WORKDIR /data
COPY go.* ./
RUN go mod download
COPY . .
RUN go test ./... -vet=off && \
	go build -ldflags="-s -w -extldflags=-static"

FROM alpine:3.18.2
WORKDIR /
COPY --from=go-app /data/mt-hosting-manager /.
CMD ["/mt-hosting-manager"]
