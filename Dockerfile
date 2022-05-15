FROM golang:1.18 as builder

WORKDIR /app
COPY . ./
RUN go mod download
RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o uru -trimpath -ldflags="-w -s" .

FROM alpine:latest as osslsigncode

RUN /bin/sh -c set -ex && apk -U upgrade && apk add libstdc++ curl ca-certificates bash tar automake autoconf libtool libcurl curl-dev libressl-dev autoconf g++ make && \
curl -SsLo osslsigncode.tar.gz "https://github.com/mtrojnar/osslsigncode/releases/download/2.3/osslsigncode-2.3.0.tar.gz" && tar -xvf osslsigncode.tar.gz && \
cd osslsigncode-2.3.0 &&  ./configure && make && make install && make clean

FROM alpine:latest

RUN apk add go openssl libtool libcurl && addgroup -S app && adduser -S app -G app
COPY --from=osslsigncode /usr/local/bin/osslsigncode /usr/local/bin/osslsigncode
USER app
WORKDIR /app
COPY --from=builder /app/uru . 
ENV PATH="/usr/lib/go/bin:${PATH}"
ENV GOROOT="/usr/lib/go/"
RUN go get github.com/C-Sto/BananaPhone
RUN go install mvdan.cc/garble@latest

ENTRYPOINT ["./uru"]