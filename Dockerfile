# build stage
FROM golang AS build
ENV GO111MODULE=on
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/lurker ./cmd/lurker

# final stage
FROM alpine:3.10
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app /go
COPY --from=build /go/src/app/config config
EXPOSE 8888
ENTRYPOINT /go/bin/lurker
