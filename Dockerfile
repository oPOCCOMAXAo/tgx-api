FROM alpine:3.19 as tdlib-builder

ENV LANG=en_US.UTF-8
ENV TZ=UTC

ARG TD_COMMIT=d7203eb719304866a7eb7033ef03d421459335b8

RUN apk update
RUN apk upgrade
RUN apk add --update \
    build-base \
    ca-certificates \
    ccache \
    cmake \
    git \
    gperf \
    linux-headers \
    openssl-dev \
    php \
    php-ctype \
    readline-dev \
    zlib-dev

RUN git clone "https://github.com/tdlib/td.git" /src

WORKDIR /src
RUN git checkout ${TD_COMMIT}
RUN mkdir ./build

WORKDIR /src/build
RUN cmake \
    -DCMAKE_BUILD_TYPE=Release \
    -DCMAKE_INSTALL_PREFIX:PATH=/usr/local \
    ..
RUN cmake --build . --target prepare_cross_compiling

WORKDIR /src
RUN php SplitSource.php

WORKDIR /src/build
RUN cmake --build . --target install -j4

RUN ls -lah /usr/local

FROM golang:1.22-alpine3.19 as go-builder

ENV LANG=en_US.UTF-8
ENV TZ=UTC

RUN set -eux 
RUN apk update
RUN apk upgrade
RUN apk add \
    bash \
    build-base \
    ca-certificates \
    curl \
    git \
    linux-headers \
    openssl-dev \
    zlib-dev

COPY --from=tdlib-builder /usr/local/include/td /usr/local/include/td/
COPY --from=tdlib-builder /usr/local/lib/libtd* /usr/local/lib/

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build \
    -a \
    -trimpath \
    -ldflags "-s -w" \
    -o /bin/app \
    "cmd/app/main.go"

FROM alpine:3.19

ENV LANG=en_US.UTF-8
ENV TZ=UTC
ENV SERVER_PORT=8080
ENV GIN_MODE=release

RUN apk upgrade --no-cache
RUN apk add --no-cache \
    ca-certificates \
    libstdc++

COPY --from=go-builder /bin/app /bin/app

EXPOSE 8080
ENTRYPOINT ["/bin/app"]

