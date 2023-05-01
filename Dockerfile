FROM golang:1.19 as builder

WORKDIR /app
COPY . ./
RUN go mod download
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o uru -trimpath -ldflags="-w -s" .

FROM alpine:latest as osslsigncode

RUN /bin/sh -c set -ex && apk -U upgrade && apk add libstdc++ cmake curl ca-certificates bash tar automake autoconf libtool libcurl curl-dev libressl-dev autoconf g++ make && \
curl -SsLo osslsigncode.tar.gz "https://github.com/mtrojnar/osslsigncode/releases/download/2.5/osslsigncode-2.5.tar.gz" && tar -xvf osslsigncode.tar.gz && \
cd osslsigncode-2.5 && mkdir build && cd build && cmake -S .. &&  cmake --build . && cmake --install .

FROM alpine:latest

RUN apk add go openssl git libtool libcurl && addgroup -S app && adduser -S app -G app
COPY --from=osslsigncode /usr/local/bin/osslsigncode /usr/local/bin/osslsigncode
RUN mkdir -p /app
RUN chown app /app
USER app
WORKDIR /app
COPY --from=builder /app/uru . 
ENV PATH="/usr/lib/go/bin:${PATH}"
ENV GOROOT="/usr/lib/go/"
RUN go install mvdan.cc/garble@latest

ENTRYPOINT ["./uru"]