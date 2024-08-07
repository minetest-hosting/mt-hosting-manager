FROM node:22.6.0-alpine as node-app
WORKDIR /public
COPY /public/package-lock.json /public/package.json ./
RUN npm ci
COPY public/ .
RUN npm run jshint && \
	npm run bundle

FROM golang:1.22.4 as go-app
WORKDIR /data
COPY go.* ./
RUN go mod download
COPY . .
COPY --from=node-app /public /data/public
RUN go test ./... && \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build .

FROM alpine:3.20.2
COPY --from=go-app /data/mt-hosting-manager /.
EXPOSE 8080
CMD ["/mt-hosting-manager"]
