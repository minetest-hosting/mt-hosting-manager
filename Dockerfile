FROM node:20.5.1-alpine as node-app
WORKDIR /public
COPY /public/package-lock.json /public/package.json ./
RUN npm ci
COPY public/ .
RUN npm run jshint && \
	npm run bundle

FROM golang:1.20.7 as go-app
WORKDIR /data
COPY go.* ./
RUN go mod download
COPY . .
COPY --from=node-app /public /data/public
RUN go test ./... -vet=off && \
	go build -ldflags="-s -w -extldflags=-static"

FROM alpine:3.18.2
WORKDIR /
COPY --from=go-app /data/mt-hosting-manager /.
CMD ["/mt-hosting-manager"]
