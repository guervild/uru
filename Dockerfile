FROM golang:1.17 as builder

WORKDIR /app
COPY . ./
RUN go mod download
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o uru -ldflags="-w -s" .

FROM alpine:latest

RUN apk add go && addgroup -S app && adduser -S app -G app
USER app
WORKDIR /app
COPY --from=builder /app/uru . 
ENV PATH="/usr/lib/go/bin:${PATH}"
ENV GOROOT="/usr/lib/go/"
ENTRYPOINT ["./uru"]